package commands

import (
    "fmt"
    "github.com/execut/ozon-parser/ozon-search-position-parser-service/internal/service/keyword"
    "os"
)

type Commander struct {
    service *keyword.Service
}

type Handler interface {
    Handle()
}

func NewCommander(service *keyword.Service) *Commander {
    return &Commander{service}
}

func (c *Commander) Handle(commandName string) {
    switch commandName {
    case "keywordsList":
        c.KeywordsList()
    case "addKeyword":
        c.AddKeyword(os.Args[len(os.Args)-1])
    default:
        panic(fmt.Sprintf("Unexpected command %s", commandName))
    }
}
