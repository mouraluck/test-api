package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testConfig() serverConfig {
	return serverConfig{
		apiKey:             "test-key",
		rateLimitPerSecond: 10,
		rateLimitBurst:     0,
	}
}

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-API-Key", "test-key")
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	newMux(testConfig()).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := "{\"message\":\"pong\",\"service\":\"test-api\"}\n"
	if got := w.Body.String(); got != expected {
		t.Fatalf("expected body %q, got %q", expected, got)
	}

	if got := w.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("expected content-type %q, got %q", "application/json; charset=utf-8", got)
	}
}

func TestPingHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	req.Header.Set("X-API-Key", "test-key")
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	newMux(testConfig()).ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestUnauthorizedWithoutAPIKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	newMux(testConfig()).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRateLimitExceeded(t *testing.T) {
	cfg := serverConfig{
		apiKey:             "test-key",
		rateLimitPerSecond: 1,
		rateLimitBurst:     0,
	}
	handler := newMux(cfg)

	first := httptest.NewRequest(http.MethodGet, "/ping", nil)
	first.Header.Set("X-API-Key", "test-key")
	first.RemoteAddr = "127.0.0.1:12345"
	firstW := httptest.NewRecorder()
	handler.ServeHTTP(firstW, first)

	second := httptest.NewRequest(http.MethodGet, "/ping", nil)
	second.Header.Set("X-API-Key", "test-key")
	second.RemoteAddr = "127.0.0.1:12345"
	secondW := httptest.NewRecorder()
	handler.ServeHTTP(secondW, second)

	if firstW.Code != http.StatusOK {
		t.Fatalf("expected first request status %d, got %d", http.StatusOK, firstW.Code)
	}

	if secondW.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request status %d, got %d", http.StatusTooManyRequests, secondW.Code)
	}
}
