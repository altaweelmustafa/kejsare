package search

import (
	"fmt"
	"github.com/altaweelmustafa/kejsare/engine/eval"
	"time"
)

// searches with increasing depth until time runs out
func (s *Searcher) IterativeSearch(maxDepth int, timeLimitMs int) Result {
	var best Result
	deadline := time.Now().Add(time.Duration(timeLimitMs) * time.Millisecond)

	for depth := 1; depth <= maxDepth; depth++ {
		// stop if we're out of time
		if time.Now().After(deadline) {
			break
		}

		result := s.Search(depth)
		best = result

		fmt.Printf("depth %d: %s score=%d nodes=%d\n",
			depth, result.Move, result.Score, result.Nodes)

		// stop early if we found checkmate
		if result.Score >= eval.Checkmate-maxDepth || result.Score <= -eval.Checkmate+maxDepth {
			break
		}
	}

	return best
}
