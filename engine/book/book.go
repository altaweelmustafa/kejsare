package book

import "math/rand"

// book maps position FEN (pieces only) to a slice of good moves
// Multiple moves per position adds variety
var book = map[string][]string{
	// Starting position
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR": {"e2e4", "d2d4", "g1f3", "c2c4"},

	// After 1.e4
	"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR": {"e7e5", "c7c5", "e7e6", "c7c6"},

	// After 1.d4
	"rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR": {"d7d5", "g8f6", "f7f5", "e7e6"},

	// After 1.c4
	"rnbqkbnr/pppppppp/8/8/2P5/8/PP1PPPPP/RNBQKBNR": {"e7e5", "c7c5", "g8f6"},

	// After 1.Nf3
	"rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R": {"d7d5", "g8f6", "c7c5"},

	// After 1.e4 e5
	"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR": {"g1f3", "f1c4", "f2f4"},

	// After 1.e4 c5 (Sicilian)
	"rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR": {"g1f3", "b1c3", "f2f4"},

	// After 1.e4 e6 (French)
	"rnbqkbnr/pppp1ppp/4p3/8/4P3/8/PPPP1PPP/RNBQKBNR": {"d2d4", "d2d3"},

	// After 1.e4 c6 (Caro-Kann)
	"rnbqkbnr/pp1ppppp/2p5/8/4P3/8/PPPP1PPP/RNBQKBNR": {"d2d4", "g1f3"},

	// After 1.e4 d6 (Pirc)
	"rnbqkbnr/ppp1pppp/3p4/8/4P3/8/PPPP1PPP/RNBQKBNR": {"d2d4", "g1f3"},

	// After 1.d4 d5
	"rnbqkbnr/ppp1pppp/8/3p4/3P4/8/PPP1PPPP/RNBQKBNR": {"c2c4", "g1f3", "b1c3"},

	// After 1.d4 Nf6 (Indian)
	"rnbqkb1r/pppppppp/5n2/8/3P4/8/PPP1PPPP/RNBQKBNR": {"c2c4", "g1f3", "b1c3"},

	// After 1.e4 e5 2.Nf3
	"rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R": {"b8c6", "g8f6", "d7d6"},

	// After 1.e4 e5 2.Nf3 Nc6
	"r1bqkbnr/pppp1ppp/2n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R": {"f1b5", "f1c4", "d2d4"},

	// After 1.e4 e5 2.Nf3 Nc6 3.Bc4 (Italian)
	"r1bqkbnr/pppp1ppp/2n5/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R": {"f8c5", "g8f6", "f7f5"},

	// After 1.d4 d5 2.c4 (Queen's Gambit)
	"rnbqkbnr/ppp1pppp/8/3p4/2PP4/8/PP2PPPP/RNBQKBNR": {"e7e6", "c7c6", "d5c4"},

	// After 1.d4 d5 2.c4 e6 (QGD)
	"rnbqkbnr/ppp2ppp/4p3/3p4/2PP4/8/PP2PPPP/RNBQKBNR": {"b1c3", "g1f3"},

	// After 1.d4 Nf6 2.c4 (King's Indian setup)
	"rnbqkb1r/pppppppp/5n2/8/2PP4/8/PP2PPPP/RNBQKBNR": {"g7g6", "e7e6", "c7c5"},

	// ── Depth 3-4 (responses to responses) ───────────────────────────────────────

	// After 1.e4 e5 2.Nf3 Nc6 3.Bc4 Bc5 (Italian Game)
	"r1bqk1nr/pppp1ppp/2n5/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQK2R": {"c2c3", "b2b4", "d2d3"},

	// After 1.e4 e5 2.Nf3 Nc6 3.Bb5 a6 (Ruy Lopez)
	"r1bqkbnr/1ppp1ppp/p1n5/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R": {"f1b5", "b1c3"},

	// After 1.e4 e5 2.Nf3 Nc6 3.Bb5 a6 4.Ba4 (Ruy Lopez main line)
	"r1bqkbnr/1ppp1ppp/p1n5/4p3/B3P3/5N2/PPPP1PPP/RNBQK2R": {"g8f6", "d7d6", "b7b5"},

	// After 1.e4 c5 2.Nf3 d6 (Sicilian Najdorf setup)
	"rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R": {"d7d6", "b8c6", "e7e6"},

	// After 1.e4 c5 2.Nf3 d6 3.d4 (Sicilian main)
	"rnbqkbnr/pp2pppp/3p4/2p5/3PP3/5N2/PPP2PPP/RNBQKB1R": {"c5d4", "g8f6"},

	// After 1.d4 d5 2.c4 e6 3.Nc3 Nf6 (QGD)
	"rnbqkb1r/ppp2ppp/4pn2/3p4/2PP4/2N5/PP2PPPP/R1BQKBNR": {"c1g5", "g1f3", "e2e3"},

	// After 1.d4 Nf6 2.c4 g6 (King's Indian)
	"rnbqkb1r/pppppp1p/5np1/8/2PP4/8/PP2PPPP/RNBQKBNR": {"b1c3", "g1f3", "e2e4"},

	// After 1.d4 Nf6 2.c4 g6 3.Nc3 Bg7 (King's Indian main)
	"rnbqk2r/ppppppbp/5np1/8/2PP4/2N5/PP2PPPP/R1BQKBNR": {"e2e4", "g1f3"},

	// After 1.e4 e5 2.Nf3 Nc6 3.Bc4 Nf6 (Two Knights)
	"r1bqkb1r/pppp1ppp/2n2n2/4p3/2B1P3/5N2/PPPP1PPP/RNBQK2R": {"d2d3", "b1c3", "g1g5"},

	// After 1.e4 e6 2.d4 d5 (French main)
	"rnbqkbnr/ppp2ppp/4p3/3p4/3PP3/8/PPP2PPP/RNBQKBNR": {"b1c3", "b1d2", "e4e5"},

	// After 1.e4 c6 2.d4 d5 (Caro-Kann main)
	"rnbqkbnr/pp2pppp/2p5/3p4/3PP3/8/PPP2PPP/RNBQKBNR": {"b1c3", "b1d2", "e4e5"},

	// After 1.d4 d5 2.c4 c6 (Slav Defense)
	"rnbqkbnr/pp2pppp/2p5/3p4/2PP4/8/PP2PPPP/RNBQKBNR": {"g1f3", "b1c3", "e2e3"},

	// After 1.e4 e5 2.f4 (King's Gambit)
	"rnbqkbnr/pppp1ppp/8/4p3/4PP2/8/PPPP2PP/RNBQKBNR": {"e5f4", "d7d5", "f8c5"},

	// After 1.Nf3 d5 2.d4 (transposing to Queen's pawn)
	"rnbqkbnr/ppp1pppp/8/3p4/3P4/5N2/PPP1PPPP/RNBQKB1R": {"g8f6", "c7c6", "e7e6"},

	// After 1.e4 e5 2.Nf3 Nc6 3.d4 (Scotch Game)
	"r1bqkbnr/pppp1ppp/2n5/4p3/3PP3/5N2/PPP2PPP/RNBQKB1R": {"e5d4", "d7d6"},

	// After 1.e4 e5 2.Nf3 Nc6 3.d4 exd4 (Scotch main)
	"r1bqkbnr/pppp1ppp/2n5/8/3pP3/5N2/PPP2PPP/RNBQKB1R": {"f3d4", "c1g5"},

	// After 1.c4 e5 2.Nc3 (English)
	"rnbqkbnr/pppp1ppp/8/4p3/2P5/2N5/PP1PPPPP/R1BQKBNR": {"g8f6", "b8c6", "f8b4"},
}

// Lookup returns a random book move for the position, or "" if not in book
func Lookup(fen string) string {
	// use only the piece placement field
	position := ""
	for _, ch := range fen {
		if ch == ' ' {
			break
		}
		position += string(ch)
	}

	moves, ok := book[position]
	if !ok || len(moves) == 0 {
		return ""
	}

	return moves[rand.Intn(len(moves))]
}
