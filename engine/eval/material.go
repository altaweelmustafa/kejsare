package eval

import (
    "github.com/altaweelmustafa/kejsare/engine/board"
)

// Piece values in centipawns 
const (
    PawnValue   = 100
    KnightValue = 320
    BishopValue = 330
    RookValue   = 500
    QueenValue  = 900
    KingValue   = 20000
)

var pieceValues = [6]int{
    PawnValue,
    KnightValue,
    BishopValue,
    RookValue,
    QueenValue,
    KingValue,
}

// Material returns the raw material score from white perspective
func Material(b *board.Board) int {
    score := 0
    for pt := board.Pawn; pt <= board.King; pt++ {
        whiteCount := b.Pieces[board.White][pt].Count()
        blackCount := b.Pieces[board.Black][pt].Count()
        score += (whiteCount - blackCount) * pieceValues[pt]
    }
    return score
}
