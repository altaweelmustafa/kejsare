package search

import (
	"github.com/altaweelmustafa/kejsare/engine/board"
	"github.com/altaweelmustafa/kejsare/engine/eval"
	"github.com/altaweelmustafa/kejsare/engine/moves"
	"github.com/altaweelmustafa/kejsare/engine/types"
)

// holds the best move and its score
type Result struct {
	Move  types.Move
	Score int
	Depth int
	Nodes int
}

// holds search state
type Searcher struct {
	board *board.Board
	nodes int
}

// creates a new Searcher
func NewSearcher(b *board.Board) *Searcher {
	return &Searcher{board: b}
}

// finds the best move at the given depth
func (s *Searcher) Search(depth int) Result {
	s.nodes = 0

	score, move := s.alphaBeta(depth, -eval.Infinity, eval.Infinity)

	return Result{
		Move:  move,
		Score: score,
		Depth: depth,
		Nodes: s.nodes,
	}
}

// returns (score, bestMove) for the current position
func (s *Searcher) alphaBeta(depth, alpha, beta int) (int, types.Move) {
	s.nodes++

	// base case, evaluate the position at depth 0
	if depth == 0 {
		return eval.Evaluate(s.board), types.NullMove
	}

	gen := moves.NewGenerator(s.board)
	ml := gen.Generate()

	// no legal moves: checkmate or stalemate
	if len(ml) == 0 {
		if s.board.IsInCheck(s.board.SideToMove) {
			// checkmate, worst possible score (adjusted by depth
			// so the engine prefers faster checkmates)
			if s.board.SideToMove == board.White {
				return -eval.Checkmate + depth, types.NullMove
			}
			return eval.Checkmate - depth, types.NullMove
		}
		// stalemate
		return 0, types.NullMove
	}

	bestMove := ml[0]
	bestScore := -eval.Infinity

	if s.board.SideToMove == board.Black {
		bestScore = eval.Infinity
	}

	for _, m := range ml {
		snap := s.board.MakeMove(m)

		score, _ := s.alphaBeta(depth-1, alpha, beta)

		s.board.UnmakeMove(m, snap)

		if s.board.SideToMove == board.White {
			// white maximizes
			if score > bestScore {
				bestScore = score
				bestMove = m
			}
			if score > alpha {
				alpha = score
			}
		} else {
			// black minimizes
			if score < bestScore {
				bestScore = score
				bestMove = m
			}
			if score < beta {
				beta = score
			}
		}

		// prune: this branch won't be chosen
		if alpha >= beta {
			break
		}
	}

	return bestScore, bestMove
}
