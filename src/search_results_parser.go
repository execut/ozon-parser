package main

import (
    "context"
    "domain"
    "encoding/json"
    "execut/ozon_parser/token/chromeCookie"
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

func OnlySimpleProducts(result domain.AnalyticsData) {
    for _, v := range result.Items {
        if v.IsTraforetto || v.SearchPromotionEnabled {
            continue
        }

        fmt.Println(fmt.Sprintf("https://www.ozon.ru/product/%s at position %v (Promo: %t)", v.Sku, v.Position, v.IsInPromo))
        reviewsPages := ParseReviews(v.Sku)
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

func CalculateAverageRankForQuery(result domain.AnalyticsData) {
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

func ParseAnalyticsForQuery(keyword domain.Keyword) domain.AnalyticsData {
    query := keyword.Name
    query = strings.TrimSpace(query)
    dataJson := SendAnalyticsHttpRequest(query)

    var result domain.AnalyticsData

    dataJson = strings.ReplaceAll(dataJson, ":\"NaN\"", ":null")
    err := json.Unmarshal([]byte(dataJson), &result)
    if err != nil {
        panic(err)
    }

    return result
}

var ctx context.Context = nil

func getChromeDpContext() context.Context {
    if ctx != nil {
        return ctx
    }

    config := cu.NewConfig(
        cu.WithHeadless(),
        cu.WithTimeout(60*time.Second),
    )
    ctx, _, err := cu.New(config)
    if err != nil {
        panic(err)
    }

    return ctx
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

func SendAnalyticsHttpRequest(query string) string {
    queryEscaped := url.QueryEscape(query)
    url := "https://seller.ozon.ru/api/explainer-service/v1/explanation?companyId=201236&query=" + queryEscaped + "&locationUid=0c5b2444-70a0-4932-980c-b4dc0d3f02b5&limit=108"

    val, err := GetCachedValue(url)
    if err != redis.Nil {
        return val
    }

    errorMessage := ""
    for errorsCount := 0; errorsCount < 2; errorsCount++ {
        var ctx = getChromeDpContext()
        var dataJson string
        var err = chromedp.Run(ctx,
            setCookie("__Secure-access-token", GetToken(), ".ozon.ru", "/", true, true),
            setCookie("__Secure-refresh-token", GetToken(), ".ozon.ru", "/", true, true),
            chromedp.Navigate(url),
            chromedp.InnerHTML(`pre`, &dataJson),
        )

        if err != nil {
            panic(err)
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

    panic(errorMessage)
}

func GetToken() string {
    token, _ := chromeCookie.ReadToken()

    return token
}
