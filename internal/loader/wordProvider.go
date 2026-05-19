package loader

import (
	"encoding/json"
	"math/rand"
	"os"
)

type WordsData struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
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
