package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
    limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
    if err != nil {
        limit = 10 // default limit
    }

    window, err := time.ParseDuration(os.Getenv("RATE_LIMIT_WINDOW"))
    if err != nil {
        window = time.Minute // default window
    }

    // Use Redis-based rate limiter
    rl := NewRedisRateLimiter(limit, window)

    mux := http.NewServeMux()
    mux.Handle("/", RedisRateLimitMiddleware(rl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })))

    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", mux)
}