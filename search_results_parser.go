package main

import (
	"context"
	"domain"
	"encoding/json"
	"execut/ozon_parser/token/chromeCookie"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"math"
	"net/url"
	"strings"
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

var ctx = context.Background()

func SendAnalyticsHttpRequest(query string) string {
	queryEscaped := url.QueryEscape(query)
	url := "https://seller.ozon.ru/api/explainer-service/v1/explanation?companyId=201236&query=" + queryEscaped + "&locationUid=0c5b2444-70a0-4932-980c-b4dc0d3f02b5&limit=108"

	val, err := GetCachedValue(url)
	if err != redis.Nil {
		return val
	}

	driver := GetSeleniumDriver()

	errorMessage := ""
	for errorsCount := 0; errorsCount < 2; errorsCount++ {
		err = driver.Get(url)
		if err != nil {
			panic(err)
		}

		element, err := driver.FindElement(selenium.ByCSSSelector, "pre")
		if err != nil {
			panic(err)
		}

		dataJson, err := element.Text()
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

var driver selenium.WebDriver

func GetSeleniumDriver() selenium.WebDriver {
	if driver != nil {
		return driver
	}

	// Run Chrome browser
	_, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		panic(err)
	}

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		// "--headless",  // comment out this line to see the browser
	},
		Path: "/usr/local/etc/chrome-114/chrome"})

	driver, err = selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}

	err = driver.Get("https://seller.ozon.ru/api/explainer-service/v1/explanation")
	if err != nil {
		panic(err)
	}

	token := GetToken()

	cookie := &selenium.Cookie{
		Name:   "__Secure-access-token",
		Value:  token,
		Expiry: math.MaxUint32,
	}
	err = driver.AddCookie(cookie)

	if err != nil {
		panic(err)
	}

	return driver
}

func GetToken() string {
	token, _ := chromeCookie.ReadToken()

	return token
}
