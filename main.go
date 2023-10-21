package main

import (
	"domain"
	"encoding/csv"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
)

func main() {

	//parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
	//positionsCmd := flag.NewFlagSet("positions", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'reviews' or 'positions' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "positions":
		keywords := extractWordsFromFile("words-for-positions.csv")
		parseKeywords(keywords)

		f, err := os.Create("positions.csv")
		defer f.Close()
		if err != nil {

			log.Fatalln("failed to open file", err)
		}

		w := csv.NewWriter(f)
		w.Comma = ';'
		defer w.Flush()

		positions := ExtractPositionsList(keywords, true)
		header := []string{
			"Date",
			"Keyword",
			"Position",
			"SKU",
			"IsInPromo",
			"IsTraforetto",
			"IsInTraforettoCampaign",
			"SearchPromotionEnabled",
			"SearchPromotionBoost",
			"PopularityScore",
			"PopularityTotalScore",
			"RatingScore",
			"Rating",
			"RatesCount",
			"PriceScore",
			"PriceDefectRateBoost",
			"PriceRub",
			"SalesScore",
			"QueryFitScore",
			"FinalResult",
		}
		w.Write(header)
		for _, position := range positions {
			var item domain.AnalyticsItem
			for _, item = range position.Data.Items {
				if !item.IsCurSellerItem {
					continue
				}

				record := []string{
					position.Time.Format("2006-01-02"),
					position.Keyword.Name,
					fmt.Sprintf("%v", item.Position),
				}

				record = append(record,
					fmt.Sprintf("%s", item.Sku),
					fmt.Sprintf("%t", item.IsInPromo),
					fmt.Sprintf("%t", item.IsTraforetto),
					fmt.Sprintf("%t", item.IsInTraforettoCampaign),
					fmt.Sprintf("%t", item.SearchPromotionEnabled),
					fmt.Sprintf("%v", item.SearchPromotionBoost),
					fmt.Sprintf("%v", item.PopularityScore),
					fmt.Sprintf("%v", item.PopularityTotalScore),
					fmt.Sprintf("%v", item.RatingScore),
					fmt.Sprintf("%v", item.Rating),
					fmt.Sprintf("%v", item.RatesCount),
					fmt.Sprintf("%v", item.PriceScore),
					fmt.Sprintf("%v", item.PriceDefectRateBoost),
					fmt.Sprintf("%s", item.PriceRub),
					fmt.Sprintf("%v", item.SalesScore),
					fmt.Sprintf("%v", item.QueryFitScore),
					fmt.Sprintf("%v", item.FinalResult),
				)

				if err := w.Write(record); err != nil {
					log.Fatalln("error writing record to file", err)
				}
			}
		}
	case "scores":
		keywords := extractWordsFromFile("words-for-scores.csv")
		positions := ExtractPositionsList(keywords, true)
		fmt.Println("Current competitors ratings:")
		fmt.Println("IsCurSellerItem;IsTrafaret;Position;PopularityScore;PopularityTotalScore;RatingScore;PriceScore;SalesScore;FinalResult")
		var targetItem domain.AnalyticsItem = domain.AnalyticsItem{}
		for _, position := range positions {
			for _, item := range position.Data.Items {
				if !item.IsCurSellerItem && (item.IsTraforetto || item.SearchPromotionEnabled || item.IsInPromo) {
					continue
				}

				fmt.Println(fmt.Sprintf("%t;%t;%v;%v;%v;%v;%v;%v;%v", item.IsCurSellerItem, item.IsTraforetto, item.Position, item.PopularityScore, item.PopularityTotalScore, item.RatingScore, item.PriceScore, item.SalesScore, item.FinalResult))
				if !item.IsCurSellerItem {
					targetItem = item
				}
			}
		}

		positions = ExtractPositionsList(keywords, false)
		fmt.Println("Rating stats:")
		fmt.Println("Date;Position;SKU;IsTrafaret;PopularityScore;PopularityTotalScore;RatingScore;PriceScore;SalesScore;FinalResult")
		rowFormat := "%s;%v;%s;%t;%v;%v;%v;%v;%v;%v"
		for _, position := range positions {
			for _, item := range position.Data.Items {
				if !item.IsCurSellerItem {
					continue
				}

				fmt.Println(fmt.Sprintf(rowFormat, position.Time.Format("2006-01-02"), item.Position, item.Sku, item.IsTraforetto, item.PopularityScore, item.PopularityTotalScore, item.RatingScore, item.PriceScore, item.SalesScore, item.FinalResult))
			}
		}

		fmt.Println(fmt.Sprintf(rowFormat, "Target:", targetItem.Position, targetItem.Sku, targetItem.IsTraforetto, targetItem.PopularityScore, targetItem.PopularityTotalScore, targetItem.RatingScore, targetItem.PriceScore, targetItem.SalesScore, targetItem.FinalResult))
	case "best-positions":
		keywords := extractWordsFromFile("words-for-best-positions.csv")
		parseKeywords(keywords)
		positions := ExtractPositionsList(keywords, true)
		fmt.Println("Link;SKU;Position")
		for _, position := range positions {
			for _, item := range position.Data.Items {
				if item.IsTraforetto || item.SearchPromotionEnabled || item.IsInPromo {
					continue
				}

				fmt.Println(fmt.Sprintf("https://www.ozon.ru/context/detail/id/%s ;%v;%v", item.Sku, item.Sku, item.Position))
			}
		}
	case "reviews":
		keywords := extractWordsFromFile("words-for-reviews.csv")
		parseKeywords(keywords)
		positions := ExtractPositionsList(keywords, true)
		for _, position := range positions {
			result := ExtractWordsFromReviews(position.Data)
			fmt.Println("Reviews word frequency for " + position.Keyword.Name + ":")
			fmt.Println("Positive")
			for _, word := range result.PositiveWords {
				if word.Value < 25 {
					break
				}

				fmt.Println(fmt.Sprintf("%s;%v", word.Key, word.Value))
			}

			fmt.Println("Negative: ")
			for _, word := range result.NegativeWords {
				if word.Value < 10 {
					break
				}

				fmt.Println(fmt.Sprintf("%s;%v", word.Key, word.Value))
			}
		}
	default:
		fmt.Println("expected 'parse' or 'positions' subcommands")
		os.Exit(1)
	}
}

func parseKeywords(keywords []domain.Keyword) {
	bar := progressbar.Default(int64(len(keywords)))
	for _, keyword := range keywords {
		bar.Describe(fmt.Sprintf("Begin parse products for %s", keyword.Name))
		result := ParseAnalyticsForQuery(keyword)
		SaveAnalyticsForQuery(keyword, result)
		e := bar.Add(1)
		if e != nil {
			panic(e)
		}
	}
}

func extractWordsFromFile(fileName string) []domain.Keyword {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	var keywords []domain.Keyword
	for {
		records, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		for _, query := range records {
			keyword := domain.Keyword{Name: query}
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}
