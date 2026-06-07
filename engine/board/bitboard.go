package board

import "math/bits" // Easier dealing with bits (ofc)

// returns the index of the least significant bit
func (bb Bitboard) LSB() int {
    return bits.TrailingZeros64(uint64(bb))
}

// clears the least significant bit and returns its index
// iterates over all set squares in a bitboard
func (bb *Bitboard) PopLSB() int {
    sq := bb.LSB()
    *bb &= *bb - 1 // clears the lowest set bit
    return sq
}

// returns the number of set bits (pieces on the board)
func (bb Bitboard) Count() int {
    return bits.OnesCount64(uint64(bb))
}

// returns true if no bits are set
func (bb Bitboard) IsEmpty() bool {
    return bb == 0
}

// returns true if any bit is set
func (bb Bitboard) NotEmpty() bool {
    return bb != 0
}
