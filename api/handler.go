package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/ws", wsHandler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
