package api

// what the frontend sends
type MoveRequest struct {
    FEN      string `json:"fen"`
    Depth    int    `json:"depth"`
    TimeLimitMs int `json:"time_limit_ms"`
}

// what the engine sends back
type MoveResponse struct {
    Move  string `json:"move"`
    Score int    `json:"score"`
    Depth int    `json:"depth"`
    Nodes int    `json:"nodes"`
    Error string `json:"error,omitempty"`
}
