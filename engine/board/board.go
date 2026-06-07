package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/altaweelmustafa/kejsare/engine/types"
)

// Color represents whose turn it is
type Color uint8

const (
	White Color = iota // iota = 0
	Black              // iota = 1
)

// Piece types — we store color and type separately
type PieceType uint8

const (
	Pawn   PieceType = iota // 0
	Knight                  // 1
	Bishop                  // 2
	Rook                    // 3
	Queen                   // 4
	King                    // 5
)

// Snapshot saves the irreversible state before a move
// so UnmakeMove can fully restore the position
type Snapshot struct {
	CastlingRights uint8
	EnPassantSq    int
	HalfMoveClock  int
	CapturedPiece  int // -1 if no capture, else PieceType index
}

// Bitboard is just a uint64 alias — makes code more readable
type Bitboard uint64

// File masks
const (
	FileA Bitboard = 0x0101010101010101 << 0
	FileB Bitboard = 0x0101010101010101 << 1
	FileC Bitboard = 0x0101010101010101 << 2
	FileD Bitboard = 0x0101010101010101 << 3
	FileE Bitboard = 0x0101010101010101 << 4
	FileF Bitboard = 0x0101010101010101 << 5
	FileG Bitboard = 0x0101010101010101 << 6
	FileH Bitboard = 0x0101010101010101 << 7
)

// Rank masks
const (
	Rank1 Bitboard = 0xFF << (8 * 0)
	Rank2 Bitboard = 0xFF << (8 * 1)
	Rank3 Bitboard = 0xFF << (8 * 2)
	Rank4 Bitboard = 0xFF << (8 * 3)
	Rank5 Bitboard = 0xFF << (8 * 4)
	Rank6 Bitboard = 0xFF << (8 * 5)
	Rank7 Bitboard = 0xFF << (8 * 6)
	Rank8 Bitboard = 0xFF << (8 * 7)
)

// Board holds the complete game state
type Board struct {
	// One bitboard per piece type per color (12 total)
	Pieces [2][6]Bitboard // [Color][PieceType]

	// Occupancy shortcuts (union of all piece bitboards)
	Occupied  Bitboard // all pieces
	OccupiedW Bitboard // white pieces only
	OccupiedB Bitboard // black pieces only

	// Game state
	SideToMove     Color
	CastlingRights uint8 // 4 bits: KQkq
	EnPassantSq    int   // -1 if none, else square index 0-63
	HalfMoveClock  int   // for 50-move rule
	FullMoveNumber int
}

// Square converts file (0-7) and rank (0-7) to square index (0-63)
func Square(file, rank int) int {
	return rank*8 + file
}

// File extracts file (0=a, 7=h) from square index
func File(sq int) int {
	return sq % 8
}

// Rank extracts rank (0=rank1, 7=rank8) from square index
func Rank(sq int) int {
	return sq / 8
}

// SquareBB returns a Bitboard with only the given square set
func SquareBB(sq int) Bitboard {
	return Bitboard(1) << sq
}

// Set sets a bit at square sq
func (bb *Bitboard) Set(sq int) {
	*bb |= Bitboard(1) << sq
}

// Clear clears a bit at square sq
func (bb *Bitboard) Clear(sq int) {
	*bb &^= Bitboard(1) << sq
}

// IsSet returns true if bit at sq is 1
func (bb Bitboard) IsSet(sq int) bool {
	return bb&(Bitboard(1)<<sq) != 0
}

// updateOccupancy rebuilds the occupancy bitboards from piece bitboards
// Call this after any move
func (b *Board) updateOccupancy() {
	b.OccupiedW = 0
	b.OccupiedB = 0
	for pt := Pawn; pt <= King; pt++ {
		b.OccupiedW |= b.Pieces[White][pt]
		b.OccupiedB |= b.Pieces[Black][pt]
	}
	b.Occupied = b.OccupiedW | b.OccupiedB
}

// NewBoard returns an empty board (no pieces)
func NewBoard() *Board {
	return &Board{
		EnPassantSq:    -1,
		FullMoveNumber: 1,
	}
}

// Maps FEN character to (Color, PieceType)
var fenPiece = map[rune]struct {
	color     Color
	pieceType PieceType
}{
	'P': {White, Pawn},
	'N': {White, Knight},
	'B': {White, Bishop},
	'R': {White, Rook},
	'Q': {White, Queen},
	'K': {White, King},
	'p': {Black, Pawn},
	'n': {Black, Knight},
	'b': {Black, Bishop},
	'r': {Black, Rook},
	'q': {Black, Queen},
	'k': {Black, King},
}

// parses a FEN string and returns a Board
func ParseFEN(fen string) (*Board, error) {
	b := NewBoard()

	// split into 6 fields
	fields := strings.Split(fen, " ")
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid FEN: expected 6 fields, got %d", len(fields))
	}

	// Field 1: piece placement
	ranks := strings.Split(fields[0], "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("invalid FEN: expected 8 ranks, got %d", len(ranks))
	}

	// ranks[0] = rank 8, ranks[7] = rank 1
	for rankIdx, rankStr := range ranks {
		rank := 7 - rankIdx // convert: rankIdx 0 -> rank 7 (rank 8 on board)
		file := 0

		for _, ch := range rankStr {
			if ch >= '1' && ch <= '8' {
				// skip that # of empty squares
				file += int(ch - '0')
			} else {
				piece, ok := fenPiece[ch]
				if !ok {
					return nil, fmt.Errorf("invalid FEN: unknown piece '%c'", ch)
				}
				sq := Square(file, rank)
				b.Pieces[piece.color][piece.pieceType].Set(sq)
				file++
			}
		}

		if file != 8 {
			return nil, fmt.Errorf("invalid FEN: rank '%s' has %d files", rankStr, file)
		}
	}

	// Field 2: turn
	switch fields[1] {
	case "w":
		b.SideToMove = White
	case "b":
		b.SideToMove = Black
	default:
		return nil, fmt.Errorf("invalid FEN: side to move must be 'w' or 'b'")
	}

	// Field 3: castling rights
	// We encode as 4 bits: b0=K, b1=Q, b2=k, b3=q
	b.CastlingRights = 0
	if fields[2] != "-" {
		for _, ch := range fields[2] {
			switch ch {
			case 'K':
				b.CastlingRights |= 1 << 0
			case 'Q':
				b.CastlingRights |= 1 << 1
			case 'k':
				b.CastlingRights |= 1 << 2
			case 'q':
				b.CastlingRights |= 1 << 3
			default:
				return nil, fmt.Errorf("invalid FEN: unknown castling right '%c'", ch)
			}
		}
	}

	// Field 4: en passant square
	b.EnPassantSq = -1
	if fields[3] != "-" {
		if len(fields[3]) != 2 {
			return nil, fmt.Errorf("invalid FEN: bad en passant square '%s'", fields[3])
		}
		epFile := int(fields[3][0] - 'a') // a=0, h=7
		epRank := int(fields[3][1] - '1') // 1=0, 8=7
		b.EnPassantSq = Square(epFile, epRank)
	}

	// Field 5: 50-rule timer
	hmc, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("invalid FEN: bad half move clock '%s'", fields[4])
	}
	b.HalfMoveClock = hmc

	// Field 6: full move number
	fmn, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, fmt.Errorf("invalid FEN: bad full move number '%s'", fields[5])
	}
	b.FullMoveNumber = fmn

	// rebuild occupancy bitboards from piece bitboards
	b.updateOccupancy()

	return b, nil
}

// +++++++++ USING AI HELP ++++++++++
// pieceChar maps (Color, PieceType) back to a display character
var pieceChar = [2][6]rune{
	{'P', 'N', 'B', 'R', 'Q', 'K'}, // White
	{'p', 'n', 'b', 'r', 'q', 'k'}, // Black
}

// String prints the board visually, rank 8 at top
func (b *Board) String() string {
	var sb strings.Builder

	sb.WriteString("\n  a b c d e f g h\n")
	sb.WriteString("  ----------------\n")

	for rank := 7; rank >= 0; rank-- {
		sb.WriteString(fmt.Sprintf("%d|", rank+1))

		for file := 0; file < 8; file++ {
			sq := Square(file, rank)
			written := false

			for color := White; color <= Black; color++ {
				for pt := Pawn; pt <= King; pt++ {
					if b.Pieces[color][pt].IsSet(sq) {
						sb.WriteRune(pieceChar[color][pt])
						sb.WriteRune(' ')
						written = true
					}
				}
			}

			if !written {
				sb.WriteString(". ")
			}
		}

		sb.WriteString(fmt.Sprintf("|%d\n", rank+1))
	}

	sb.WriteString("  ----------------\n")
	sb.WriteString("  a b c d e f g h\n")

	return sb.String()
}

// MakeMove applies a move to the board, mutating it in place
// Returns a snapshot you can pass to UnmakeMove to restore the position
func (b *Board) MakeMove(m types.Move) Snapshot {
	// save state we can't reconstruct from the move alone
	snap := Snapshot{
		CastlingRights: b.CastlingRights,
		EnPassantSq:    b.EnPassantSq,
		HalfMoveClock:  b.HalfMoveClock,
		CapturedPiece:  -1,
	}

	from := m.From()
	to := m.To()
	flags := m.Flags()
	us := b.SideToMove
	them := Black
	if us == Black {
		them = White
	}

	// find which piece is moving
	movingPiece := b.pieceAt(us, from)

	// handle capture
	if m.IsCapture() && !m.IsEnPassant() {
		captured := b.pieceAt(them, to)
		if captured != -1 {
			b.Pieces[them][captured].Clear(to)
			snap.CapturedPiece = captured
		}
	}

	// move the piece
	b.Pieces[us][movingPiece].Clear(from)
	b.Pieces[us][movingPiece].Set(to)

	// special cases
	switch {
	case flags == types.FlagEnPassant:
		// captured pawn is behind the 'to' square
		capSq := to - 8
		if us == Black {
			capSq = to + 8
		}
		b.Pieces[them][Pawn].Clear(capSq)
		snap.CapturedPiece = int(Pawn)

	case flags == types.FlagKingCastle:
		// move rook kingside
		if us == White {
			b.Pieces[White][Rook].Clear(7) // h1
			b.Pieces[White][Rook].Set(5)   // f1
		} else {
			b.Pieces[Black][Rook].Clear(63) // h8
			b.Pieces[Black][Rook].Set(61)   // f8
		}

	case flags == types.FlagQueenCastle:
		// move rook queenside
		if us == White {
			b.Pieces[White][Rook].Clear(0) // a1
			b.Pieces[White][Rook].Set(3)   // d1
		} else {
			b.Pieces[Black][Rook].Clear(56) // a8
			b.Pieces[Black][Rook].Set(59)   // d8
		}

	case m.IsPromotion():
		// replace pawn with promoted piece
		b.Pieces[us][Pawn].Clear(to)
		promoType := promotionPiece(flags)
		b.Pieces[us][promoType].Set(to)
	}

	// update castling rights
	// any king or rook move revokes the relevant right
	b.CastlingRights &= castlingMask[from] & castlingMask[to]

	// update en passant square
	b.EnPassantSq = -1
	if flags == types.FlagDoublePush {
		// ep square is the square the pawn skipped over
		if us == White {
			b.EnPassantSq = to - 8
		} else {
			b.EnPassantSq = to + 8
		}
	}

	// update clocks
	if m.IsCapture() || movingPiece == int(Pawn) {
		b.HalfMoveClock = 0
	} else {
		b.HalfMoveClock++
	}
	if us == Black {
		b.FullMoveNumber++
	}

	// flip side to move
	b.SideToMove = them

	b.updateOccupancy()

	return snap
}

// pieceAt returns the PieceType index for the piece of color c on square sq
// returns -1 if no piece found
func (b *Board) pieceAt(c Color, sq int) int {
	for pt := int(Pawn); pt <= int(King); pt++ {
		if b.Pieces[c][pt].IsSet(sq) {
			return pt
		}
	}
	return -1
}

// promotionPiece maps promotion flags to PieceType
func promotionPiece(flags int) PieceType {
	switch flags {
	case types.FlagKnightPromo, types.FlagKnightPromoCap:
		return Knight
	case types.FlagBishopPromo, types.FlagBishopPromoCap:
		return Bishop
	case types.FlagRookPromo, types.FlagRookPromoCap:
		return Rook
	default:
		return Queen
	}
}

// castlingMask revokes castling rights when a piece moves from/to a key square
// e.g. if the h1 rook moves, white loses kingside castling
var castlingMask = [64]uint8{
	13, 15, 15, 15, 12, 15, 15, 14, // rank 1 (a1=13 revokes Qa1, h1=14 revokes Ka1)
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	7, 15, 15, 15, 3, 15, 15, 11, // rank 8 (a8=7 revokes qa8, h8=11 revokes ka8)
}

// UnmakeMove reverses a move using the saved Snapshot
func (b *Board) UnmakeMove(m types.Move, snap Snapshot) {
	from := m.From()
	to := m.To()
	flags := m.Flags()

	// flip side back — 'us' is who made the move
	b.SideToMove ^= 1
	us := b.SideToMove
	them := Black
	if us == Black {
		them = White
	}

	// find the moving piece (after promotion it's the promo piece)
	movingPiece := b.pieceAt(us, to)

	// reverse promotion first
	if m.IsPromotion() {
		b.Pieces[us][movingPiece].Clear(to)
		b.Pieces[us][Pawn].Set(to)
		movingPiece = int(Pawn)
	}

	// move piece back
	b.Pieces[us][movingPiece].Clear(to)
	b.Pieces[us][movingPiece].Set(from)

	// restore captured piece
	if m.IsCapture() && !m.IsEnPassant() && snap.CapturedPiece != -1 {
		b.Pieces[them][snap.CapturedPiece].Set(to)
	}

	// reverse special cases
	switch {
	case flags == types.FlagEnPassant:
		capSq := to - 8
		if us == Black {
			capSq = to + 8
		}
		b.Pieces[them][Pawn].Set(capSq)

	case flags == types.FlagKingCastle:
		if us == White {
			b.Pieces[White][Rook].Clear(5)
			b.Pieces[White][Rook].Set(7)
		} else {
			b.Pieces[Black][Rook].Clear(61)
			b.Pieces[Black][Rook].Set(63)
		}

	case flags == types.FlagQueenCastle:
		if us == White {
			b.Pieces[White][Rook].Clear(3)
			b.Pieces[White][Rook].Set(0)
		} else {
			b.Pieces[Black][Rook].Clear(59)
			b.Pieces[Black][Rook].Set(56)
		}
	}

	// restore saved state
	b.CastlingRights = snap.CastlingRights
	b.EnPassantSq = snap.EnPassantSq
	b.HalfMoveClock = snap.HalfMoveClock
	if us == Black {
		b.FullMoveNumber--
	}

	b.updateOccupancy()
}

// returns true if the given color's king is attacked
func (b *Board) IsInCheck(c Color) bool {
    kingBB := b.Pieces[c][King]
    if kingBB.IsEmpty() {
        return false
    }
    kingSq := kingBB.LSB()
    return b.IsAttacked(kingSq, c)
}

// returns true if square sq is attacked by any enemy of color us
func (b *Board) IsAttacked(sq int, us Color) bool {
    them := Black
    if us == Black {
        them = White
    }

    // check knight attacks
    if KnightMask(sq)&b.Pieces[them][Knight] != 0 {
        return true
    }

    // check pawn attacks
    if PawnAttackMask(sq, us)&b.Pieces[them][Pawn] != 0 {
        return true
    }

    // check king attacks
    if KingMask(sq)&b.Pieces[them][King] != 0 {
        return true
    }

    // check diagonal (bishop/queen)
    if BishopAttacks(sq, b.Occupied)&(b.Pieces[them][Bishop]|b.Pieces[them][Queen]) != 0 {
        return true
    }

    // check orthogonal (rook/queen)
    if RookAttacks(sq, b.Occupied)&(b.Pieces[them][Rook]|b.Pieces[them][Queen]) != 0 {
        return true
    }

    return false
}
