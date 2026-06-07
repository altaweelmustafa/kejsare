package board


func BishopAttacks(sq int, occupied Bitboard) Bitboard {
    return rayAttacks(sq, occupied,  9, FileA) |
           rayAttacks(sq, occupied,  7, FileH) |
           rayAttacks(sq, occupied, -9, FileH) |
           rayAttacks(sq, occupied, -7, FileA)
}

func RookAttacks(sq int, occupied Bitboard) Bitboard {
    return rayAttacks(sq, occupied,  8, 0)    |
           rayAttacks(sq, occupied, -8, 0)    |
           rayAttacks(sq, occupied,  1, FileA) |
           rayAttacks(sq, occupied, -1, FileH)
}

func rayAttacks(sq int, occupied Bitboard, shift int, fileMask Bitboard) Bitboard {
    attacks := Bitboard(0)
    bb      := SquareBB(sq)

    for {
        if shift > 0 {
            bb = (bb << shift) & ^fileMask
        } else {
            bb = (bb >> (-shift)) & ^fileMask
        }
        if bb == 0 {
            break
        }
        attacks |= bb
        if bb&occupied != 0 {
            break
        }
    }
    return attacks
}

func KnightMask(sq int) Bitboard {
    bb := SquareBB(sq)
    return (bb<<17)&^FileA |
           (bb<<15)&^FileH |
           (bb<<10)&^(FileA|FileB) |
           (bb<<6) &^(FileG|FileH) |
           (bb>>17)&^FileH |
           (bb>>15)&^FileA |
           (bb>>10)&^(FileG|FileH) |
           (bb>>6) &^(FileA|FileB)
}

func KingMask(sq int) Bitboard {
    bb := SquareBB(sq)
    return (bb<<8) |
           (bb>>8) |
           ((bb<<1) & ^FileA) |
           ((bb>>1) & ^FileH) |
           ((bb<<9) & ^FileA) |
           ((bb<<7) & ^FileH) |
           ((bb>>9) & ^FileH) |
           ((bb>>7) & ^FileA)
}

func PawnAttackMask(sq int, us Color) Bitboard {
    bb := SquareBB(sq)
    if us == White {
        // squares a white king on sq would be attacked from by black pawns
        return ((bb << 7) & ^FileH) | ((bb << 9) & ^FileA)
    }
    return ((bb >> 7) & ^FileA) | ((bb >> 9) & ^FileH)
}
