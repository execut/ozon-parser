package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"math"
	"net/url"
	"os"
	"strings"
)

func OnlySimpleProducts(result AnaliticsData) {
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

func CalculateAverageRankForQuery(result AnaliticsData) {
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

func ParseAnalyticsForQuery(query string) AnaliticsData {
	query = strings.TrimSpace(query)
	dataJson := SendAnalyticsHttpRequest(query)

	var result AnaliticsData

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

	SetCachedValue(url, dataJson)

	return dataJson
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
	}})

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
	tokenB, err := os.ReadFile("token.txt")
	if err != nil {
		panic(err)
	}

	token := string(tokenB)

	return token
}

type AnaliticsData struct {
	Items []struct {
		Position               int     `json:"position"`
		IsTraforetto           bool    `json:"isTraforetto"`
		Sku                    string  `json:"sku"`
		Name                   string  `json:"name"`
		ImageUrl               string  `json:"imageUrl"`
		SellerName             string  `json:"sellerName"`
		IsCurSellerItem        bool    `json:"isCurSellerItem"`
		IsInPromo              bool    `json:"isInPromo"`
		Delivery               string  `json:"delivery"`
		DeliveryBoost          int     `json:"deliveryBoost"`
		PopularityScore        int     `json:"popularityScore"`
		SalesScore             int     `json:"salesScore"`
		PriceRub               string  `json:"priceRub"`
		PriceScore             int     `json:"priceScore"`
		Rating                 float64 `json:"rating"`
		RatesCount             int     `json:"ratesCount"`
		RatingScore            int     `json:"ratingScore"`
		QueryFitScore          int     `json:"queryFitScore"`
		PopularityTotalScore   int     `json:"popularityTotalScore"`
		DeliverySpeed          string  `json:"deliverySpeed"`
		FinalResult            float64 `json:"finalResult"`
		SearchPromotionBoost   float64 `json:"searchPromotionBoost"`
		SearchPromotionEnabled bool    `json:"searchPromotionEnabled"`
		PriceDefectRateBoost   float64 `json:"priceDefectRateBoost"`
		IsInTraforettoCampaign bool    `json:"isInTraforettoCampaign"`
		DeliverySpeedBoostSlot *struct {
			FromDays int `json:"fromDays"`
			ToDays   int `json:"toDays"`
		} `json:"deliverySpeedBoostSlot"`
	} `json:"items"`
	Analytics struct {
		CurSellerItems struct {
			ItemsInTopQnty int           `json:"itemsInTopQnty"`
			ItemsTotalQnty int           `json:"itemsTotalQnty"`
			PagesQueried   int           `json:"pagesQueried"`
			PagesResult    int           `json:"pagesResult"`
			ItemsOutOfTop  []interface{} `json:"itemsOutOfTop"`
		} `json:"curSellerItems"`
		ExpressDelivery struct {
			CurSellerItemsQnty      int `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int `json:"curSellerItemsTotalQnty"`
		} `json:"expressDelivery"`
		LocalStore struct {
			CurSellerItemsQnty      int `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int `json:"curSellerItemsTotalQnty"`
		} `json:"localStore"`
		Delivery struct {
			DeliveryType            string `json:"deliveryType"`
			CurSellerItemsQnty      int    `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int    `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int    `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int    `json:"curSellerItemsTotalQnty"`
		} `json:"delivery"`
		Promo struct {
			CurSellerItemsQnty      int `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty    int `json:"competitorsItemsQnty"`
			CompetitorsItemsTopSize int `json:"competitorsItemsTopSize"`
			CurSellerItemsTotalQnty int `json:"curSellerItemsTotalQnty"`
		} `json:"promo"`
		ItemsRating struct {
			CurSellerItemsAvgRating   float64 `json:"curSellerItemsAvgRating"`
			CompetitorsItemsAvgRating float64 `json:"competitorsItemsAvgRating"`
			CompetitorsItemsTopSize   int     `json:"competitorsItemsTopSize"`
			CurSellerItemsMinRating   int     `json:"curSellerItemsMinRating"`
			CompetitorsItemsMaxRating float64 `json:"competitorsItemsMaxRating"`
			CurSellerMinRatingItem    string  `json:"curSellerMinRatingItem"`
		} `json:"itemsRating"`
		SearchPromotion struct {
			CurSellerItemsQnty            int      `json:"curSellerItemsQnty"`
			CompetitorsItemsQnty          int      `json:"competitorsItemsQnty"`
			MaxCompetitorsSearchPromotion float64  `json:"maxCompetitorsSearchPromotion"`
			CurSellerNotPromotedItems     []string `json:"curSellerNotPromotedItems"`
			CurSellerItemsTotalQnty       int      `json:"curSellerItemsTotalQnty"`
		} `json:"searchPromotion"`
		QueryFit struct {
			CurSellerItemsMinScore   int    `json:"curSellerItemsMinScore"`
			CompetitorsItemsMaxScore int    `json:"competitorsItemsMaxScore"`
			CurSellerMinScoreItem    string `json:"curSellerMinScoreItem"`
		} `json:"queryFit"`
	} `json:"analytics"`
}
