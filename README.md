# typest

A simple terminal typing test built with Bubble Tea.

## What it does

- Shows a target sentence in the terminal
- Highlights correct and incorrect typed characters
- Advances word by word as you type
- Exits when the sentence is complete

## Requirements

- Go 1.26+

## Run

```bash
go run .
```

## Controls

- Type to match the displayed text
- Space moves to the next word when the current word is correct
- Ctrl+C quits

## Project files

- main.go: app model, update loop, and rendering
- styles.go: text styles for typed, untyped, cursor, and errors

## Notes

This project is currently focused on English text input and a simple typing flow.
