package main

import (
    "flag"
    "github.com/execut/ozon-parser/ozon-search-position-parser-service/app/commands"
    "github.com/execut/ozon-parser/ozon-search-position-parser-service/internal/service/keyword"
)

var action string

func main() {
    flag.Parse()
    keywordService := keyword.NewService()

    commander := commands.NewCommander(keywordService)

    commander.Handle(action)
}

func init() {
    flag.StringVar(&action, "action", "", "")
}
