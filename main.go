package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("words.csv")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	for {
		records, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		for _, query := range records {
			var reviewsText = ""
			result := ParseAnalyticsForQuery(query)
			for _, v := range result.Items {
				if v.IsTraforetto || v.SearchPromotionEnabled {
					continue
				}

				reviewsPages := ParseReviews(v.Sku)
				reviewsText += ExtractTextFromReviews(reviewsPages)
			}

			words := CountWords(reviewsText)
			for _, word := range words {
				if word.value < 5 {
					continue
				}

				fmt.Println(word.key, ": ", word.value)
			}
		}
	}
}
