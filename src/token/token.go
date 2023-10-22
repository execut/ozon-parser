package token

import (
    "github.com/browserutils/kooky"
    "github.com/browserutils/kooky/browser/chrome"
    "os"
)

type Token interface {
    Value() string
}

type Chrome struct{}

type File struct {
    FilePath string
}

func (c Chrome) Value() string {
    dir, _ := os.UserConfigDir() // "/<USER>/Library/Application Support/"
    cookiesFile := dir + "/google-chrome/Default/Cookies"
    cookies, err := chrome.ReadCookies(cookiesFile, kooky.Valid, kooky.DomainHasSuffix(`ozon.ru`), kooky.Name(`__Secure-access-token`))
    if err != nil {
        panic(err)
    }

    for _, cookie := range cookies {
        return cookie.Value
    }

    panic("Cookie not found")
}

func (f File) Value() string {
    b, _ := os.ReadFile(f.FilePath)

    return string(b)
}
