package main

import (
	"github.com/bbalet/stopwords"
	"github.com/imbue11235/words"
	"gitlab.com/opennota/morph"
	"sort"
	"strings"
)

type Word struct {
	key   string
	value int
}

func sortedWords(words map[string]int) []Word {
	var sorted []Word
	for k, v := range words {
		sorted = append(sorted, Word{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})
	return sorted
}

var morphInitialized = false

func CountWords(text string) []Word {
	text = strings.ToLower(text)
	text = stopwords.CleanString(text, "ru", false)

	wordsForNormalization := words.Extract(text)
	// loading the dictionary data;
	// you can also use morph.InitWith("path/to/the/dictionary")
	if !morphInitialized {
		if err := morph.Init(); err != nil {
			panic(err)
		}

		morphInitialized = true
	}

	normalizedWords := make(map[string]int)
	for _, wordForNormalization := range wordsForNormalization {
		// parsing
		words, norms, _ := morph.Parse(wordForNormalization)
		var currentWord string
		if len(words) > 0 {
			currentWord = norms[0]
		} else {
			currentWord = wordForNormalization
		}

		normalizedWords[currentWord]++
	}

	return sortedWords(normalizedWords)
}
