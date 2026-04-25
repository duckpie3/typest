# typest

A simple terminal typing test built with Bubble Tea.

## What it does

- Loads a random quote for each test
- Highlights correct and incorrect typed characters in real time
- Advances word by word as you type
- Tracks elapsed time and WPM while typing
- Shows a post-test results screen with WPM stats and a WPM-over-time graph
- Lets you start a new test directly from the results screen

## Requirements

- Go 1.26+

## Run

```bash
go run ./cmd/typest
```

## Controls

- Type to match the displayed text
- Space moves to the next word when the current word is correct
- Ctrl+Backspace clears the current word
- On the results screen, press Space to start another test
- Ctrl+C quits

## Project files

```
typest/
├── assets/                    ← datasets
│   ├── quotes.json
│   └── words.json
├── bin/
├── cmd/
│   └── typest/
│       └── main.go
├── go.mod
├── go.sum
├── internal/
│   ├── app/
│   │   └── app.go             ← app model, update loop, and screen transitions
│   ├── quotes/
│   │   └── loader.go          ← quote data loading and lookup helpers
│   ├── results/
│   │   └── results.go         ← results view and WPM graph rendering
│   ├── typing/
│   │   ├── stats.go
│   │   └── typing.go          ← typing test state, input handling, live stats collection
│   └── ui/
│       └── styles.go          ← text styles and test area layout
├── README.md
└── tests/
```

## Notes

- The quote dataset in `quotes.json` is sourced from the Monkeytype repository.
- This project is currently focused on English text input and a simple typing flow.
