package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/altaweelmustafa/kejsare/api"
    "github.com/altaweelmustafa/kejsare/config"
)

func main() {
    cfg := config.Load()

    mux := http.NewServeMux()
    api.RegisterRoutes(mux)

    addr := fmt.Sprintf(":%s", cfg.Port)
    log.Printf("kejsare engine listening on %s", addr)

    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}
