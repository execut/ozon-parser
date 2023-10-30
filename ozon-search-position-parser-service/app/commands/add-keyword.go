package commands

import "fmt"

func (c Commander) AddKeyword(keyword string) {
    fmt.Println(fmt.Sprintf("Keyword \"%s\" added", keyword))
}
