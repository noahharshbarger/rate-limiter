package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "requests_total",
            Help: "Total number of requests",
        },
        []string{"status"},
    )
)

func init() {
    prometheus.MustRegister(requestsTotal)
}

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

func RateLimitMiddleware(rl *RateLimiter, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := getIP(r)
        if !rl.IsAllowed(ip) {
            log.Printf("Rate limit exceeded for IP: %s", ip)
            requestsTotal.WithLabelValues("rate_limited").Inc()
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        log.Printf("Request allowed for IP: %s", ip)
        requestsTotal.WithLabelValues("allowed").Inc()
        next.ServeHTTP(w, r)
    })
}

func main() {
    limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
    if err != nil {
        limit = 10 // default limit
    }

    window, err := time.ParseDuration(os.Getenv("RATE_LIMIT_WINDOW"))
    if err != nil {
        window = time.Minute // default window
    }

    rl := NewRateLimiter(limit, window)

    mux := http.NewServeMux()
    mux.Handle("/", RateLimitMiddleware(rl, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    })))

    // Expose the Prometheus metrics endpoint
    mux.Handle("/metrics", promhttp.Handler())

    log.Println("Starting server on :8080")
    http.ListenAndServe(":8080", mux)
}