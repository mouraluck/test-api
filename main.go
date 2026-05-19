package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type serverConfig struct {
	apiKey             string
	rateLimitPerSecond int
	rateLimitBurst     int
}

type rateLimiter struct {
	mu          sync.Mutex
	windowStart time.Time
	counts      map[string]int
	limit       int
}

func newRateLimiter(limitPerSecond, burst int) *rateLimiter {
	return &rateLimiter{
		windowStart: time.Now(),
		counts:      make(map[string]int),
		limit:       limitPerSecond + burst,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.windowStart) >= time.Second {
		rl.windowStart = now
		rl.counts = make(map[string]int)
	}

	rl.counts[key]++
	return rl.counts[key] <= rl.limit
}

func loadConfigFromEnv() (serverConfig, error) {
	cfg := serverConfig{
		apiKey:             os.Getenv("API_KEY"),
		rateLimitPerSecond: 5,
		rateLimitBurst:     2,
	}

	if cfg.apiKey == "" {
		return serverConfig{}, errors.New("API_KEY is required")
	}

	if v := os.Getenv("RATE_LIMIT_PER_SECOND"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed <= 0 {
			return serverConfig{}, fmt.Errorf("invalid RATE_LIMIT_PER_SECOND: %q", v)
		}
		cfg.rateLimitPerSecond = parsed
	}

	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed < 0 {
			return serverConfig{}, fmt.Errorf("invalid RATE_LIMIT_BURST: %q", v)
		}
		cfg.rateLimitBurst = parsed
	}

	return cfg, nil
}

func newMux(cfg serverConfig) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)

	limiter := newRateLimiter(cfg.rateLimitPerSecond, cfg.rateLimitBurst)
	return withAPIKeyAuth(cfg.apiKey, withRateLimit(limiter, mux))
}

func withAPIKeyAuth(expectedKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != expectedKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withRateLimit(limiter *rateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientKey := clientIP(r.RemoteAddr)
		if !limiter.allow(clientKey) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
		"service": "test-api",
	})
}

func main() {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		log.Fatalf("invalid server config: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, newMux(cfg)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
