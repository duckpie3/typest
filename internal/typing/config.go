package typing

import "strconv"

type Mode int

const (
	ModeQuotes Mode = iota
	ModeWords
)

func (m Mode) String() string {
	switch m {
	case ModeQuotes:
		return "quotes"
	case ModeWords:
		return "words"
	default:
		return "unknown"
	}
}

type QuoteLength int

const (
	QuoteSmall QuoteLength = iota
	QuoteMedium
	QuoteLong
)

func (q QuoteLength) String() string {
	switch q {
	case QuoteSmall:
		return "small"
	case QuoteMedium:
		return "medium"
	case QuoteLong:
		return "long"
	default:
		return "unknown"
	}
}

const (
	WordsShort  = 10
	WordsMedium = 25
	WordsLong   = 50
)

type Config struct {
	Mode        Mode
	QuoteLength QuoteLength
	WordsCount  int
}

func DefaultConfig() Config {
	return Config{Mode: ModeQuotes, QuoteLength: QuoteMedium, WordsCount: WordsMedium}
}

func ToggleMode(cfg Config) Config {
	if cfg.Mode == ModeQuotes {
		cfg.Mode = ModeWords
		return cfg
	}
	cfg.Mode = ModeQuotes
	return cfg
}

func ApplyLengthPreset(cfg Config, preset int) Config {
	switch cfg.Mode {
	case ModeQuotes:
		switch preset {
		case 1:
			cfg.QuoteLength = QuoteSmall
		case 2:
			cfg.QuoteLength = QuoteMedium
		case 3:
			cfg.QuoteLength = QuoteLong
		}
	case ModeWords:
		switch preset {
		case 1:
			cfg.WordsCount = WordsShort
		case 2:
			cfg.WordsCount = WordsMedium
		case 3:
			cfg.WordsCount = WordsLong
		}
	}
	return cfg
}

func (cfg Config) LengthLabel() string {
	if cfg.Mode == ModeQuotes {
		return cfg.QuoteLength.String()
	}
	return strconv.Itoa(cfg.WordsCount)
}
