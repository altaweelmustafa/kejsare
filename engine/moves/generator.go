package moves

import (
	"github.com/altaweelmustafa/kejsare/engine/board"
	"github.com/altaweelmustafa/kejsare/engine/types"
)

// holds a reference to the board and builds move lists
type Generator struct {
	b *board.Board
}

// creates a Generator for the given board
func NewGenerator(b *board.Board) *Generator {
	return &Generator{b: b}
}

// returns all legal moves for the side to move
func (g *Generator) Generate() []types.Move {
	moves := make([]types.Move, 0, 40) // 40 is typical move count per position

	if g.b.SideToMove == board.White {
		moves = g.generatePawns(moves, board.White)
		moves = g.generateKnights(moves, board.White)
		moves = g.generateBishops(moves, board.White)
		moves = g.generateRooks(moves, board.White)
		moves = g.generateQueens(moves, board.White)
		moves = g.generateKing(moves, board.White)
	} else {
		moves = g.generatePawns(moves, board.Black)
		moves = g.generateKnights(moves, board.Black)
		moves = g.generateBishops(moves, board.Black)
		moves = g.generateRooks(moves, board.Black)
		moves = g.generateQueens(moves, board.Black)
		moves = g.generateKing(moves, board.Black)
	}

	return FilterLegal(g.b, moves)
}

//  ============= Sliding helpers ================

// rayAttacks casts a ray from sq in a given direction until blocked
// direction is encoded as a shift amount (positive = left shift, negative = right)
// wrap guards prevent the ray from wrapping across the board edge
func rayAttacks(sq int, occupied board.Bitboard, shift int, fileMask board.Bitboard) board.Bitboard {
	attacks := board.Bitboard(0)
	bb := board.SquareBB(sq)

	for {
		if shift > 0 {
			bb = (bb << shift) & ^fileMask
		} else {
			bb = (bb >> (-shift)) & ^fileMask
		}

		if bb == 0 {
			break // hit the board edge
		}

		attacks |= bb

		if bb&occupied != 0 {
			break // hit a piece — include the square but stop
		}
	}

	return attacks
}

// bishopAttacks returns all squares a bishop on sq can attack
func bishopAttacks(sq int, occupied board.Bitboard) board.Bitboard {
	return rayAttacks(sq, occupied, 9, board.FileA) | // northeast
		rayAttacks(sq, occupied, 7, board.FileH) | // northwest
		rayAttacks(sq, occupied, -9, board.FileH) | // southwest
		rayAttacks(sq, occupied, -7, board.FileA) // southeast
}

// rookAttacks returns all squares a rook on sq can attack
func rookAttacks(sq int, occupied board.Bitboard) board.Bitboard {
	return rayAttacks(sq, occupied, 8, 0) | // north
		rayAttacks(sq, occupied, -8, 0) | // south
		rayAttacks(sq, occupied, 1, board.FileA) | // east
		rayAttacks(sq, occupied, -1, board.FileH) // west
}

// ===== Pawns ==========

func (g *Generator) generatePawns(ml []types.Move, c board.Color) []types.Move {
	pawns := g.b.Pieces[c][board.Pawn]
	empty := ^g.b.Occupied // squares with no pieces

	var (
		enemy     board.Bitboard
		forward   board.Bitboard
		startRank board.Bitboard
		promoRank board.Bitboard
		push      func(board.Bitboard) board.Bitboard
		pushTwo   func(board.Bitboard) board.Bitboard
		attLeft   func(board.Bitboard) board.Bitboard
		attRight  func(board.Bitboard) board.Bitboard
		fromPush  func(int) int
		fromLeft  func(int) int
		fromRight func(int) int
	)

	if c == board.White {
		enemy = g.b.OccupiedB
		startRank = board.Rank2
		promoRank = board.Rank8
		push = func(bb board.Bitboard) board.Bitboard { return bb << 8 }
		pushTwo = func(bb board.Bitboard) board.Bitboard { return bb << 16 }
		attLeft = func(bb board.Bitboard) board.Bitboard { return (bb << 7) & ^board.FileH }
		attRight = func(bb board.Bitboard) board.Bitboard { return (bb << 9) & ^board.FileA }
		fromPush = func(sq int) int { return sq - 8 }
		fromLeft = func(sq int) int { return sq - 7 }
		fromRight = func(sq int) int { return sq - 9 }
	} else {
		enemy = g.b.OccupiedW
		startRank = board.Rank7
		promoRank = board.Rank1
		push = func(bb board.Bitboard) board.Bitboard { return bb >> 8 }
		pushTwo = func(bb board.Bitboard) board.Bitboard { return bb >> 16 }
		attLeft = func(bb board.Bitboard) board.Bitboard { return (bb >> 7) & ^board.FileA }
		attRight = func(bb board.Bitboard) board.Bitboard { return (bb >> 9) & ^board.FileH }
		fromPush = func(sq int) int { return sq + 8 }
		fromLeft = func(sq int) int { return sq + 7 }
		fromRight = func(sq int) int { return sq + 9 }
	}

	_ = forward // suppress unused warning until we use it

	// Single push
	singlePush := push(pawns) & empty
	nonPromo := singlePush & ^promoRank
	promos := singlePush & promoRank

	targets := nonPromo
	for targets.NotEmpty() {
		to := targets.PopLSB()
		ml = append(ml, types.NewMove(fromPush(to), to, types.FlagQuiet))
	}

	// promotions from single push
	targets = promos
	for targets.NotEmpty() {
		to := targets.PopLSB()
		from := fromPush(to)
		ml = append(ml, types.NewMove(from, to, types.FlagQueenPromo))
		ml = append(ml, types.NewMove(from, to, types.FlagRookPromo))
		ml = append(ml, types.NewMove(from, to, types.FlagBishopPromo))
		ml = append(ml, types.NewMove(from, to, types.FlagKnightPromo))
	}

	// Double push
	doublePush := pushTwo(pawns&startRank) & empty & push(empty)
	targets = doublePush
	for targets.NotEmpty() {
		to := targets.PopLSB()
		ml = append(ml, types.NewMove(fromPush(fromPush(to)), to, types.FlagDoublePush))
	}

	// Captures
	leftCaps := attLeft(pawns) & enemy
	targets = leftCaps
	for targets.NotEmpty() {
		to := targets.PopLSB()
		from := fromLeft(to)
		if board.SquareBB(to)&promoRank != 0 {
			ml = append(ml, types.NewMove(from, to, types.FlagQueenPromoCap))
			ml = append(ml, types.NewMove(from, to, types.FlagRookPromoCap))
			ml = append(ml, types.NewMove(from, to, types.FlagBishopPromoCap))
			ml = append(ml, types.NewMove(from, to, types.FlagKnightPromoCap))
		} else {
			ml = append(ml, types.NewMove(from, to, types.FlagCapture))
		}
	}

	rightCaps := attRight(pawns) & enemy
	targets = rightCaps
	for targets.NotEmpty() {
		to := targets.PopLSB()
		from := fromRight(to)
		if board.SquareBB(to)&promoRank != 0 {
			ml = append(ml, types.NewMove(from, to, types.FlagQueenPromoCap))
			ml = append(ml, types.NewMove(from, to, types.FlagRookPromoCap))
			ml = append(ml, types.NewMove(from, to, types.FlagBishopPromoCap))
			ml = append(ml, types.NewMove(from, to, types.FlagKnightPromoCap))
		} else {
			ml = append(ml, types.NewMove(from, to, types.FlagCapture))
		}
	}

	// En passant
	if g.b.EnPassantSq != -1 {
		epBB := board.SquareBB(g.b.EnPassantSq)
		leftEP := attLeft(pawns) & epBB
		rightEP := attRight(pawns) & epBB

		if leftEP.NotEmpty() {
			to := leftEP.PopLSB()
			ml = append(ml, types.NewMove(fromLeft(to), to, types.FlagEnPassant))
		}
		if rightEP.NotEmpty() {
			to := rightEP.PopLSB()
			ml = append(ml, types.NewMove(fromRight(to), to, types.FlagEnPassant))
		}
	}

	return ml
}

// Knights

func (g *Generator) generateKnights(ml []types.Move, c board.Color) []types.Move {
	knights := g.b.Pieces[c][board.Knight]
	own := g.b.OccupiedW
	if c == board.Black {
		own = g.b.OccupiedB
	}

	for knights.NotEmpty() {
		from := knights.PopLSB()
		attacks := knightAttacks(board.SquareBB(from)) & ^own

		for attacks.NotEmpty() {
			to := attacks.PopLSB()
			if g.b.Occupied.IsSet(to) {
				ml = append(ml, types.NewMove(from, to, types.FlagCapture))
			} else {
				ml = append(ml, types.NewMove(from, to, types.FlagQuiet))
			}
		}
	}

	return ml
}

// knightAttacks returns all squares a knight on bb can jump to
func knightAttacks(bb board.Bitboard) board.Bitboard {
	return (bb<<17)&^board.FileA |
		(bb<<15)&^board.FileH |
		(bb<<10)&^(board.FileA|board.FileB) |
		(bb<<6)&^(board.FileG|board.FileH) |
		(bb>>17)&^board.FileH |
		(bb>>15)&^board.FileA |
		(bb>>10)&^(board.FileG|board.FileH) |
		(bb>>6)&^(board.FileA|board.FileB)
}

// ============== Bishops ==================

func (g *Generator) generateBishops(ml []types.Move, c board.Color) []types.Move {
	bishops := g.b.Pieces[c][board.Bishop]
	own := g.b.OccupiedW
	if c == board.Black {
		own = g.b.OccupiedB
	}

	for bishops.NotEmpty() {
		from := bishops.PopLSB()
		attacks := bishopAttacks(from, g.b.Occupied) & ^own

		for attacks.NotEmpty() {
			to := attacks.PopLSB()
			if g.b.Occupied.IsSet(to) {
				ml = append(ml, types.NewMove(from, to, types.FlagCapture))
			} else {
				ml = append(ml, types.NewMove(from, to, types.FlagQuiet))
			}
		}
	}

	return ml
}

// =============== Rooks ==================

func (g *Generator) generateRooks(ml []types.Move, c board.Color) []types.Move {
	rooks := g.b.Pieces[c][board.Rook]
	own := g.b.OccupiedW
	if c == board.Black {
		own = g.b.OccupiedB
	}

	for rooks.NotEmpty() {
		from := rooks.PopLSB()
		attacks := rookAttacks(from, g.b.Occupied) & ^own

		for attacks.NotEmpty() {
			to := attacks.PopLSB()
			if g.b.Occupied.IsSet(to) {
				ml = append(ml, types.NewMove(from, to, types.FlagCapture))
			} else {
				ml = append(ml, types.NewMove(from, to, types.FlagQuiet))
			}
		}
	}

	return ml
}

// ============== Queens ===================

func (g *Generator) generateQueens(ml []types.Move, c board.Color) []types.Move {
	queens := g.b.Pieces[c][board.Queen]
	own := g.b.OccupiedW
	if c == board.Black {
		own = g.b.OccupiedB
	}

	for queens.NotEmpty() {
		from := queens.PopLSB()
		// queen = bishop + rook combined
		attacks := (bishopAttacks(from, g.b.Occupied) |
			rookAttacks(from, g.b.Occupied)) & ^own

		for attacks.NotEmpty() {
			to := attacks.PopLSB()
			if g.b.Occupied.IsSet(to) {
				ml = append(ml, types.NewMove(from, to, types.FlagCapture))
			} else {
				ml = append(ml, types.NewMove(from, to, types.FlagQuiet))
			}
		}
	}

	return ml
}

// ================ King ==================

func (g *Generator) generateKing(ml []types.Move, c board.Color) []types.Move {
	king := g.b.Pieces[c][board.King]
	own := g.b.OccupiedW
	if c == board.Black {
		own = g.b.OccupiedB
	}

	if king.IsEmpty() {
		return ml
	}

	from := king.LSB()
	attacks := kingAttacks(king) & ^own

	for attacks.NotEmpty() {
		to := attacks.PopLSB()
		if g.b.Occupied.IsSet(to) {
			ml = append(ml, types.NewMove(from, to, types.FlagCapture))
		} else {
			ml = append(ml, types.NewMove(from, to, types.FlagQuiet))
		}
	}

	// Castling
	if c == board.White {
		// kingside: bits f1, g1 must be empty
		if g.b.CastlingRights&(1<<0) != 0 &&
			!g.b.Occupied.IsSet(5) &&
			!g.b.Occupied.IsSet(6) {
			ml = append(ml, types.NewMove(4, 6, types.FlagKingCastle))
		}
		// queenside: bits b1, c1, d1 must be empty
		if g.b.CastlingRights&(1<<1) != 0 &&
			!g.b.Occupied.IsSet(1) &&
			!g.b.Occupied.IsSet(2) &&
			!g.b.Occupied.IsSet(3) {
			ml = append(ml, types.NewMove(4, 2, types.FlagQueenCastle))
		}
	} else {
		// kingside: f8, g8 must be empty
		if g.b.CastlingRights&(1<<2) != 0 &&
			!g.b.Occupied.IsSet(61) &&
			!g.b.Occupied.IsSet(62) {
			ml = append(ml, types.NewMove(60, 62, types.FlagKingCastle))
		}
		// queenside: b8, c8, d8 must be empty
		if g.b.CastlingRights&(1<<3) != 0 &&
			!g.b.Occupied.IsSet(57) &&
			!g.b.Occupied.IsSet(58) &&
			!g.b.Occupied.IsSet(59) {
			ml = append(ml, types.NewMove(60, 58, types.FlagQueenCastle))
		}
	}

	return ml
}

func kingAttacks(bb board.Bitboard) board.Bitboard {
	return (bb << 8) | // north
		(bb >> 8) | // south
		((bb << 1) & ^board.FileA) | // east
		((bb >> 1) & ^board.FileH) | // west
		((bb << 9) & ^board.FileA) | // northeast
		((bb << 7) & ^board.FileH) | // northwest
		((bb >> 9) & ^board.FileH) | // southwest
		((bb >> 7) & ^board.FileA) // southeast
}
