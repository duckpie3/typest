package loader

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
)

type Quote struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Length int    `json:"length"`
	ID     int    `json:"id"`
}

type QuotesData struct {
	Language string  `json:"language"`
	Groups   [][]int `json:"groups"`
	Quotes   []Quote `json:"quotes"`
}

type WordsData struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

func LoadQuotes(path string) (*QuotesData, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var d QuotesData
	if err := json.Unmarshal(bytes, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d QuotesData) QuoteAt(i int) (Quote, bool) {
	if i < 0 || i >= len(d.Quotes) {
		return Quote{}, false
	}
	return d.Quotes[i], true
}

func (d QuotesData) RandomQuote() Quote {
	index := rand.Intn(len(d.Quotes))
	quote, ok := d.QuoteAt(index)
	if ok {
		return quote
	} else {
		log.Println("Error obtaining random quote.")
		return quote
	}
}

func LoadWords(path string) (*WordsData, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d WordsData
	if err := json.Unmarshal(bytes, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d WordsData) RandomWords(length int) []string {
	words := make([]string, length)
	for i := range length {
		index := rand.Intn(len(d.Words))
		words[i] = d.Words[index]
	}
	return words
}
