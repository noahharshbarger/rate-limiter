package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// In-memory RateLimiter
type RateLimiter struct {
    visitors    map[string]int
    lastRequest map[string]time.Time
    mu          sync.Mutex
    limit       int
    window      time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        visitors:    make(map[string]int),
        lastRequest: make(map[string]time.Time),
        limit:       limit,
        window:      window,
    }
}

func (rl *RateLimiter) IsAllowed(ip string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    count, exists := rl.visitors[ip]
    if !exists || time.Since(rl.lastRequest[ip]) > rl.window {
        rl.visitors[ip] = 1
        rl.lastRequest[ip] = time.Now()
        return true
    }

    if count < rl.limit {
        rl.visitors[ip]++
        rl.lastRequest[ip] = time.Now()
        return true
    }

    return false
}

// Redis-based RateLimiter
type RedisRateLimiter struct {
    client *redis.Client
    limit  int
    window time.Duration
}

func NewRedisRateLimiter(limit int, window time.Duration) *RedisRateLimiter {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Redis address
    })
    return &RedisRateLimiter{
        client: client,
        limit:  limit,
        window: window,
    }
}

func (rl *RedisRateLimiter) IsAllowed(ip string) bool {
    ctx := context.Background()
    count, err := rl.client.Get(ctx, ip).Int()
    if err == redis.Nil || count < rl.limit {
        rl.client.Incr(ctx, ip)
        rl.client.Expire(ctx, ip, rl.window)
        return true
    }
    return false
}

// Function to get the user ID from the request
// func getUserID(r *http.Request) string {
     // Extract user ID from request (e.g., from headers or query parameters)
//     return r.Header.Get("X-User-ID")
// }

// Middleware to apply rate limiting
func RateLimitMiddleware(rl *RateLimiter, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := getIP(r)
        if !rl.IsAllowed(ip) {
            log.Printf("Rate limit exceeded for IP: %s", ip)
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        log.Printf("Request allowed for IP: %s", ip)
        next.ServeHTTP(w, r)
    })
}

// Middleware to apply Redis-based rate limiting
func RedisRateLimitMiddleware(rl *RedisRateLimiter, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := getIP(r)
        if !rl.IsAllowed(ip) {
            log.Printf("Rate limit exceeded for IP: %s", ip)
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        log.Printf("Request allowed for IP: %s", ip)
        next.ServeHTTP(w, r)
    })
}

// Function to get the IP address from the request
func getIP(r *http.Request) string {
    // Check if the request is behind a proxy
    if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
        return ip
    }
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    return r.RemoteAddr
}