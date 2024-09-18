package main

import (
	"sync"
	"time"
)

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