package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/altaweelmustafa/kejsare/engine/board"
	"github.com/altaweelmustafa/kejsare/engine/book"
	"github.com/altaweelmustafa/kejsare/engine/search"
	"golang.org/x/net/websocket"
)

// wraps the websocket handler with permissive origin check
func wsHandler() http.Handler {
	srv := websocket.Server{
		Handshake: func(cfg *websocket.Config, r *http.Request) error {
			return nil // accept any origin
		},
		Handler: websocket.Handler(Handler),
	}
	return srv
}

// handles a single WebSocket connection
func Handler(ws *websocket.Conn) {
	log.Println("client connected")
	defer log.Println("client disconnected")

	for {
		var raw string
		if err := websocket.Message.Receive(ws, &raw); err != nil {
			break
		}

		var req MoveRequest
		if err := json.Unmarshal([]byte(raw), &req); err != nil {
			sendError(ws, "invalid request: "+err.Error())
			continue
		}

		if req.Depth == 0 {
			req.Depth = 6
		}
		if req.TimeLimitMs == 0 {
			req.TimeLimitMs = 3000
		}

		b, err := board.ParseFEN(req.FEN)
		if err != nil {
			sendError(ws, "invalid FEN: "+err.Error())
			continue
		}

		// opening book
		move := book.Lookup(req.FEN)
		score := 0
		nodes := 0
		depth := 0

		if move == "" {
			s := search.NewSearcher(b)
			result := s.IterativeSearch(req.Depth, req.TimeLimitMs)
			move = result.Move.String()
			score = result.Score
			nodes = result.Nodes
			depth = result.Depth
		}

		resp := MoveResponse{
			Move:  move,
			Score: score,
			Depth: depth,
			Nodes: nodes,
		}

		data, _ := json.Marshal(resp)
		websocket.Message.Send(ws, string(data))
	}
}

func sendError(ws *websocket.Conn, msg string) {
	resp := MoveResponse{Error: msg}
	data, _ := json.Marshal(resp)
	websocket.Message.Send(ws, string(data))
}
