package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
    // Set environment variables for testing
    os.Setenv("RATE_LIMIT", "2")
    os.Setenv("RATE_LIMIT_WINDOW", "1s")

    limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
    if err != nil {
        t.Fatalf("Failed to parse RATE_LIMIT: %v", err)
    }

    window, err := time.ParseDuration(os.Getenv("RATE_LIMIT_WINDOW"))
    if err != nil {
        t.Fatalf("Failed to parse RATE_LIMIT_WINDOW: %v", err)
    }

    rl := NewRedisRateLimiter(limit, window)

    mux := http.NewServeMux()
    mux.Handle("/", RedisRateLimitMiddleware(rl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })))

    server := httptest.NewServer(mux)
    defer server.Close()

    client := &http.Client{}

    // First request should be allowed
    resp, err := client.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }

    // Second request should be allowed
    resp, err = client.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }

    // Third request should be rate-limited
    resp, err = client.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusTooManyRequests {
        t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, resp.StatusCode)
    }

    // Wait for the rate limit window to expire
    time.Sleep(window)

    // Fourth request should be allowed
    resp, err = client.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
    }
}