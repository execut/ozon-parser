package main

import (
	"domain"
	"encoding/json"
	"fmt"
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"github.com/redis/go-redis/v9"
	"log"
	"math"
)

func ParseReviews(productId string) []domain.ReviewsData {
	var ok bool
	lastPage := 999
	var result []domain.ReviewsData
	// New creates a new context for use with chromedp. With this context
	// you can use chromedp as you normally would.
	config := cu.NewConfig(
		// Remove this if you want to see a browser window.
		cu.WithHeadless(),

		// If the webelement is not found within 10 seconds, timeout.
		//cu.WithTimeout(10 * time.Second),
	)
	ctx, _, err := cu.New(config)
	if err != nil {
		panic(err)
	}

	for page := 1; page <= lastPage; page++ {
		url := fmt.Sprintf("https://www.ozon.ru/product/%s/reviews/?page=%v", productId, page)
		dataState, err := GetCachedValue(url)
		if err == redis.Nil {
			err = chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.WaitVisible(`//*[@id="state-webListReviews-3231710-default-1"]`),
				chromedp.AttributeValue(`#state-webListReviews-3231710-default-1`, `data-state`, &dataState, &ok),
			)
			if err != nil {
				fmt.Println(url)
				log.Fatal(err)
			}

			//time.Sleep(time.Second * 10)

			SetCachedValue(url, dataState)
		}

		var reviewsData domain.ReviewsData
		err = json.Unmarshal([]byte(dataState), &reviewsData)
		if err != nil {
			panic(err)
		}

		if lastPage == 999 {
			lastPage = int(math.Ceil(float64(reviewsData.Paging.Total) / float64(reviewsData.Paging.PerPage)))
			if lastPage > 10 {
				lastPage = 10
			}
		}

		result = append(result, reviewsData)
	}

	return result
}

func ExtractTextFromReviews(reviewsPages []domain.ReviewsData) string {
	var text string
	for _, reviewsPage := range reviewsPages {
		for _, reviewsList := range reviewsPage.Reviews {
			content := reviewsList.Content
			if content.Positive == "" && content.Negative == "" && content.Comment == "" {
				continue
			}

			if content.Comment != "" {
				text += " " + content.Comment
			}

			if content.Positive != "" {
				text += " " + content.Positive
			}

			if content.Negative != "" {
				text += " " + content.Negative
			}
		}
	}

	return text
}

func ExtractWordsFromReviews(result domain.AnalyticsData) {
	var reviewsText = ""
	for _, v := range result.Items {
		if v.IsTraforetto || v.SearchPromotionEnabled || v.IsInPromo {
			continue
		}

		fmt.Println(fmt.Sprintf("Product page https://www.ozon.ru/product/%s/", v.Sku))
		fmt.Println(fmt.Sprintf("Parse reviews for https://www.ozon.ru/product/%s/reviews/", v.Sku))
		reviewsPages := ParseReviews(v.Sku)
		reviewsText += ExtractTextFromReviews(reviewsPages)
	}

	words := CountWords(reviewsText)
	for key, word := range words {
		if key > 25 {
			break
		}

		fmt.Println(word.key, ": ", word.value)
	}
}
