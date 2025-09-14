package presentation

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redisClient *redis.Client
	Requests    int
	Window      time.Duration
}

func NewRateLimiter(redisClient *redis.Client, requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		Requests:    requests,
		Window:      window,
	}
}

func (rl *RateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		key := "rate_limit:" + ip

		count, err := rl.redisClient.Incr(r.Context(), key).Result()
		if err != nil {
			log.Printf("Rate limiter error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			rl.redisClient.Expire(r.Context(), key, rl.Window)
		}

		if count > int64(rl.Requests) {
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", rl.Window.Seconds()))
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
