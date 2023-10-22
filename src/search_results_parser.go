package main

import (
    "context"
    "domain"
    "encoding/json"
    "execut/ozon_parser/token"
    "fmt"
    cu "github.com/Davincible/chromedp-undetected"
    "github.com/chromedp/cdproto/cdp"
    "github.com/chromedp/cdproto/network"
    "github.com/chromedp/chromedp"
    "github.com/redis/go-redis/v9"
    "net/url"
    "strings"
    "time"
)

type SearchResultsParser struct {
    token token.Token
}

func (p *SearchResultsParser) OnlySimpleProducts(result domain.AnalyticsData, reviewsParser ReviewsParser) {
    for _, v := range result.Items {
        if v.IsTraforetto || v.SearchPromotionEnabled {
            continue
        }

        fmt.Println(fmt.Sprintf("https://www.ozon.ru/product/%s at position %v (Promo: %t)", v.Sku, v.Position, v.IsInPromo))
        reviewsPages := reviewsParser.ParseReviews(v.Sku)
        for _, reviewsPage := range reviewsPages {
            for _, reviewsList := range reviewsPage.Reviews {
                fmt.Println(reviewsList.Author.FirstName, reviewsList.Author.LastName)
                fmt.Println(reviewsList.Content.Positive)
                fmt.Println(reviewsList.Content.Negative)
                fmt.Println(reviewsList.Content.Comment)
            }
        }
    }
}

func (p *SearchResultsParser) CalculateAverageRankForQuery(result domain.AnalyticsData) {
    var totalResult float64
    var totalNonTrafaretCount int = 0
    for _, v := range result.Items {
        if v.IsTraforetto {
            continue
        }

        totalResult += v.FinalResult
        totalNonTrafaretCount++
        if totalNonTrafaretCount == 50 {
            break
        }
    }

    if totalNonTrafaretCount != 50 && len(result.Items) == 108 {
        panic(fmt.Sprintf("Not enough entries to calculate average final result. Current count: %v of %v", totalNonTrafaretCount, len(result.Items)))
    }

    totalAverageResult := totalResult / float64(totalNonTrafaretCount)
    fmt.Println(fmt.Sprintf("%.4f", totalAverageResult))
}

func (p *SearchResultsParser) Parse(keyword domain.Keyword) domain.AnalyticsData {
    query := keyword.Name
    query = strings.TrimSpace(query)
    dataJson := p.SendAnalyticsHttpRequest(query)

    var result domain.AnalyticsData

    dataJson = strings.ReplaceAll(dataJson, ":\"NaN\"", ":null")
    err := json.Unmarshal([]byte(dataJson), &result)
    if err != nil {
        panic(err)
    }

    return result
}

func (p *SearchResultsParser) getChromeDpContext() (context.Context, context.CancelFunc) {
    var opts []chromedp.ContextOption
    //opts = append(opts,
    //    chromedp.WithLogf(log.Printf),
    //    chromedp.WithDebugf(log.Printf),
    //    chromedp.WithErrorf(log.Printf),
    //)

    config := cu.NewConfig(
        cu.WithHeadless(),
        cu.WithTimeout(60*time.Second),
        cu.WithLogLevel(10),
    )
    config.ContextOptions = opts
    ctx, cancel, err := cu.New(config)
    if err != nil {
        panic(err)
    }

    return ctx, cancel
}

func setCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
    return chromedp.ActionFunc(func(ctx context.Context) error {
        expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
        var err = network.SetCookie(name, value).
            WithExpires(&expr).
            WithDomain(domain).
            WithPath(path).
            WithHTTPOnly(httpOnly).
            WithSecure(secure).
            Do(ctx)
        if err != nil {
            return err
        }

        return nil
    })
}

func (p *SearchResultsParser) SendAnalyticsHttpRequest(query string) string {
    queryEscaped := url.QueryEscape(query)
    url := "https://seller.ozon.ru/api/explainer-service/v1/explanation?companyId=201236&query=" + queryEscaped + "&locationUid=0c5b2444-70a0-4932-980c-b4dc0d3f02b5&limit=108"

    val, err := GetCachedValue(url)
    if err != redis.Nil {
        return val
    }

    err = nil
    errorMessage := ""
    var ctx, cancel = p.getChromeDpContext()
    defer cancel()
    for errorsCount := 0; errorsCount < 3; errorsCount++ {
        var dataJson string
        err = chromedp.Run(ctx,
            setCookie("__Secure-access-token", p.token.Value(), ".ozon.ru", "/", true, true),
            setCookie("__Secure-refresh-token", p.token.Value(), ".ozon.ru", "/", true, true),
            chromedp.Navigate(url),
            chromedp.InnerHTML(`pre`, &dataJson),
        )

        if err != nil {
            fmt.Println("Try #", errorsCount, " because has error: ", err)
            continue
        }

        if strings.HasPrefix(dataJson, "{\"error\"") {
            panic("Token expire")
        }

        if strings.Contains(dataJson, "{\"code\":13,") {
            errorMessage = "Bad response: " + dataJson
            continue
        }

        if !strings.HasSuffix(dataJson, "}") {
            errorMessage = "Wrong suffix"
            continue
        }

        if !json.Valid([]byte(dataJson)) {
            errorMessage = "Wrong json"
            continue
        }

        SetCachedValue(url, dataJson)

        return dataJson
    }

    if err != nil {
        panic(err)
    }

    panic(errorMessage)
}
