package quotes

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

type Data struct {
	Language string  `json:"language"`
	Groups   [][]int `json:"groups"`
	Quotes   []Quote `json:"quotes"`
}

func LoadQuotes(path string) (*Data, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var d Data
	if err := json.Unmarshal(bytes, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *Data) QuoteAt(i int) (Quote, bool) {
	if i < 0 || i >= len(d.Quotes) {
		return Quote{}, false
	}
	return d.Quotes[i], true
}

func (d *Data) RandomQuote() Quote {
	index := rand.Intn(len(d.Quotes))
	// quote, ok := d.QuoteAt(59)
	quote, ok := d.QuoteAt(index)
	if ok {
		return quote
	} else {
		log.Println("Error obtaining random quote.")
		return quote
	}
}
