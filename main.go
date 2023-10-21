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
		defer w.Flush()

		positions := ExtractPositionsList(keywords)
		header := []string{
			"Date",
			"Keyword",
			"Position",
			"QueryFitScore",
			"Score",
			"Is promo",
			"Is trafaret",
			"Is search promotion boost",
			"PopularityScore",
			"PopularityTotalScore",
		}
		w.Write(header)
		for _, position := range positions {
			positionNumber := -1
			var item domain.AnalyticsItem
			for _, item = range position.Data.Items {
				if item.IsCurSellerItem {
					positionNumber = item.Position
					break
				}
			}

			record := []string{
				position.Time.Format("2006-01-02"),
				position.Keyword.Name,
				fmt.Sprintf("%v", positionNumber),
			}

			if positionNumber != -1 {
				record = append(record, fmt.Sprintf("%v", item.QueryFitScore), fmt.Sprintf("%v", item.FinalResult), fmt.Sprintf("%t", item.IsInPromo), fmt.Sprintf("%t", item.IsTraforetto), fmt.Sprintf("%v", item.SearchPromotionBoost), fmt.Sprintf("%v", item.PopularityScore), fmt.Sprintf("%v", item.PopularityTotalScore))
			}

			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	case "reviews":
		keywords := extractWordsFromFile("words-for-reviews.csv")
		parseKeywords(keywords)
		positions := ExtractPositionsList(keywords)
		for _, position := range positions {
			result := ExtractWordsFromReviews(position.Data)
			fmt.Println("Reviews word frequency for " + position.Keyword.Name + ":")
			fmt.Println("Positive: ")
			for _, word := range result.PositiveWords {
				if word.Value < 25 {
					break
				}

				fmt.Println(fmt.Sprintf("%s: %v", word.Key, word.Value))
			}

			fmt.Println("Negative: ")
			for _, word := range result.NegativeWords {
				if word.Value < 10 {
					break
				}

				fmt.Println(fmt.Sprintf("%s: %v", word.Key, word.Value))
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
