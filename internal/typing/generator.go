package typing

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/duckpie3/typest/internal/loader"
	"github.com/duckpie3/typest/internal/ui"
)

type generatedTest struct {
	words      []string
	wordsView  []string
	characters int
}

func generateTest(quotesData *loader.QuotesData, wordsData *loader.WordsData, config Config) (generatedTest, error) {
	switch config.Mode {
	case ModeQuotes:
		return buildQuoteTest(quotesData, config.QuoteLength)
	case ModeWords:
		return buildWordsTest(wordsData, config.WordsCount)
	default:
		return generatedTest{}, errors.New("unknown mode")
	}
}

func buildQuoteTest(quotesData *loader.QuotesData, preset QuoteLength) (generatedTest, error) {
	if quotesData == nil || len(quotesData.Quotes) == 0 {
		return generatedTest{}, errors.New("missing quotes data")
	}

	min, max := quoteLengthBounds(quotesData.Groups, preset)
	candidates := make([]loader.Quote, 0, len(quotesData.Quotes))
	for _, quote := range quotesData.Quotes {
		quoteLength := quote.Length
		if quoteLength <= 0 {
			quoteLength = len(quote.Text)
		}
		if quoteLength >= min && quoteLength <= max {
			candidates = append(candidates, quote)
		}
	}

	var selected loader.Quote
	if len(candidates) == 0 {
		selected = quotesData.RandomQuote()
	} else {
		selected = candidates[rand.Intn(len(candidates))]
	}

	rawWords := strings.Split(selected.Text, " ")
	words, wordsView := buildWordsView(rawWords)
	characters := selected.Length
	if characters <= 0 {
		characters = countCharacters(rawWords)
	}

	return generatedTest{words: words, wordsView: wordsView, characters: characters}, nil
}

func buildWordsTest(wordsData *loader.WordsData, count int) (generatedTest, error) {
	if wordsData == nil || len(wordsData.Words) == 0 {
		return generatedTest{}, errors.New("missing words data")
	}
	if count <= 0 {
		count = WordsMedium
	}

	rawWords := wordsData.RandomWords(count)
	words, wordsView := buildWordsView(rawWords)
	characters := countCharacters(rawWords)

	return generatedTest{words: words, wordsView: wordsView, characters: characters}, nil
}

func buildWordsView(rawWords []string) ([]string, []string) {
	words := make([]string, len(rawWords))
	wordsView := make([]string, len(rawWords))
	for i, word := range rawWords {
		words[i] = word + " "
		wordsView[i] = ui.UntypedStyle.Render(words[i])
	}
	return words, wordsView
}

func countCharacters(rawWords []string) int {
	total := 0
	for _, word := range rawWords {
		total += len(word)
	}
	return total
}

func quoteLengthBounds(groups [][]int, preset QuoteLength) (int, int) {
	switch preset {
	case QuoteSmall:
		if len(groups) > 0 && len(groups[0]) == 2 {
			return groups[0][0], groups[0][1]
		}
		return 0, 100
	case QuoteMedium:
		if len(groups) > 1 && len(groups[1]) == 2 {
			return groups[1][0], groups[1][1]
		}
		return 101, 300
	case QuoteLong:
		min := 301
		max := 9999
		if len(groups) > 2 && len(groups[2]) == 2 {
			min = groups[2][0]
			max = groups[2][1]
		}
		if len(groups) > 3 {
			last := groups[len(groups)-1]
			if len(last) == 2 && last[1] > max {
				max = last[1]
			}
		}
		return min, max
	default:
		return 0, 100
	}
}
