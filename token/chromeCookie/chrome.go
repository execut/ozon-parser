package chromeCookie

import (
	"errors"
	"github.com/browserutils/kooky"
	"github.com/browserutils/kooky/browser/chrome"
	"os"
)

func ReadToken() (string, error) {
	dir, _ := os.UserConfigDir() // "/<USER>/Library/Application Support/"
	cookiesFile := dir + "/google-chrome/Default/Cookies"
	cookies, err := chrome.ReadCookies(cookiesFile, kooky.Valid, kooky.DomainHasSuffix(`ozon.ru`), kooky.Name(`__Secure-access-token`))
	if err != nil {
		panic(err)
	}

	for _, cookie := range cookies {
		return cookie.Value, nil
	}

	return "", errors.New("Cookie not found")
}
