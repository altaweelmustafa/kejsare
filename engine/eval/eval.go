package eval

import (
    "github.com/altaweelmustafa/kejsare/engine/board"
)

const (
    Infinity  = 1_000_000
    Checkmate = 900_000  // score we assign to checkmate
)

// scores the position from white's perspective
func Evaluate(b *board.Board) int {
    return Material(b)
}
