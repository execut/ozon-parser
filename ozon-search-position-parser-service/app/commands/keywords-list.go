package commands

import (
    "fmt"
)

func (c Commander) KeywordsList() {
    keywords := c.service.List()
    for i := 0; i < len(keywords); i++ {
        keyword := keywords[i]
        fmt.Println(keyword.Name)
    }
}
