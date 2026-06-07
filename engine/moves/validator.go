package moves

import (
    "github.com/altaweelmustafa/kejsare/engine/board"
    "github.com/altaweelmustafa/kejsare/engine/types"
)

// removes moves that leave the king in check
func FilterLegal(b *board.Board, ml []types.Move) []types.Move {
    legal := ml[:0] // reuse the same slice, no allocation

    for _, m := range ml {
        snap := b.MakeMove(m)

        // check the side that just moved (before the flip, so ^1)
        us := b.SideToMove ^ 1
        if !b.IsInCheck(board.Color(us)) {
            legal = append(legal, m)
        }

        b.UnmakeMove(m, snap)
    }

    return legal
}
