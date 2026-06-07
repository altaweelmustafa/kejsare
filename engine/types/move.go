package types

// Move encodes a chess move in a single uint32
//
// bits  0-5:  from square
// bits  6-11: to square
// bits 12-15: flags
type Move uint32

// NullMove represents no move (used in search)
const NullMove Move = 0

// Flag constants
const (
	FlagQuiet          = 0b0000
	FlagDoublePush     = 0b0001
	FlagKingCastle     = 0b0010
	FlagQueenCastle    = 0b0011
	FlagCapture        = 0b0100
	FlagEnPassant      = 0b0101
	FlagKnightPromo    = 0b1000
	FlagBishopPromo    = 0b1001
	FlagRookPromo      = 0b1010
	FlagQueenPromo     = 0b1011
	FlagKnightPromoCap = 0b1100
	FlagBishopPromoCap = 0b1101
	FlagRookPromoCap   = 0b1110
	FlagQueenPromoCap  = 0b1111
)

// converts [from, to, flags] into a Move
func NewMove(from, to, flags int) Move {
	return Move(from | (to << 6) | (flags << 12))
}

// extracts the from square
func (m Move) From() int {
	return int(m) & 0x3F // 0x3F = 0b00111111 = mask bottom 6 bits
}

// extracts the to square
func (m Move) To() int {
	return (int(m) >> 6) & 0x3F // Same as From() but shift 6 bits at first
}

// extracts the flag bits
func (m Move) Flags() int {
	return (int(m) >> 12) & 0xF // shift down 12 then mask first 4 bits
}

// returns true for any capturing move
func (m Move) IsCapture() bool {
	return m.Flags()&FlagCapture != 0
}

// returns true for any promotion move
func (m Move) IsPromotion() bool {
	return m.Flags()&0b1000 != 0
}

// returns true for an en passant captures
func (m Move) IsEnPassant() bool {
	return m.Flags() == FlagEnPassant
}

// returns true for castling moves
func (m Move) IsCastle() bool {
	f := m.Flags()
	return f == FlagKingCastle || f == FlagQueenCastle
}

// String returns the move in UCI format (e.g e2e4)
func (m Move) String() string {
	files := "abcdefgh"
	from := m.From()
	to := m.To()

	result := string([]byte{
		files[from%8],
		byte(49 + from/8), // byte(49) is 1
		files[to%8],
		byte(49 + to/8),
	})

	// append promotion piece (if there is)
	switch m.Flags() {
	case FlagKnightPromo, FlagKnightPromoCap:
		result += "n"
	case FlagBishopPromo, FlagBishopPromoCap:
		result += "b"
	case FlagRookPromo, FlagRookPromoCap:
		result += "r"
	case FlagQueenPromo, FlagQueenPromoCap:
		result += "q"
	}

	return result
}
