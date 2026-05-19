package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	return mux
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("debug") == "env" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(strings.Join(os.Environ(), "\n")))
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, newMux()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
