package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {

	//parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
	//positionsCmd := flag.NewFlagSet("positions", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'parse' or 'positions' subcommands")
		os.Exit(1)
	}

	// Check which subcommand is invoked.
	switch os.Args[1] {

	// For every subcommand, we parse its own flags and
	// have access to trailing positional arguments.
	case "parse":

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
				keyword := Keyword{Name: query}
				fmt.Println(fmt.Sprintf("Begin parse products for %s", query))
				var reviewsText = ""
				result := ParseAnalyticsForQuery(query)
				SaveAnalyticsForQuery(keyword, result)
				for _, v := range result.Items {
					if v.IsTraforetto || v.SearchPromotionEnabled || v.IsInPromo {
						continue
					}

					//fmt.Println(fmt.Sprintf("Product page https://www.ozon.ru/product/%s/", v.Sku))
					//fmt.Println(fmt.Sprintf("Parse reviews for https://www.ozon.ru/product/%s/reviews/", v.Sku))
					//reviewsPages := ParseReviews(v.Sku)
					//reviewsText += ExtractTextFromReviews(reviewsPages)
				}

				words := CountWords(reviewsText)
				for key, word := range words {
					if key > 25 {
						break
					}

					fmt.Println(word.key, ": ", word.value)
				}
			}
		}
		//case "positions":
		positions := ExtractPositionsList()
		for _, position := range positions {
			for _, item := range position.Data.Items {
				if item.IsCurSellerItem {
					fmt.Println(fmt.Sprintf("%s; %s; %v", position.Time.Format("2006-01-02"), position.Keyword.Name, item.Position))
					break
				}
			}
		}
	default:
		fmt.Println("expected 'parse' or 'positions' subcommands")
		os.Exit(1)
	}
}
