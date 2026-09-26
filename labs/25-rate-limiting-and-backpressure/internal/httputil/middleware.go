package httputil

import (
	"encoding/json"
	"net/http"
	"strconv"

	"labs/25-rate-limiting-and-backpressure/internal/ratelimit"
)

type RateLimitResponse struct {
	Error      string `json:"error"`
	RetryAfter int    `json:"retry_after"`
}

// RateLimitMiddleware enforces tenant-based rate limits (RFC 6585 status 429).
func RateLimitMiddleware(registry *ratelimit.Registry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantKey := r.Header.Get("X-API-Key")
		if tenantKey == "" {
			tenantKey = "anonymous"
		}

		bucket := registry.Get(tenantKey)
		if !bucket.Allow() {
			retryAfter := bucket.RetryAfterSeconds(1.0)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			w.WriteHeader(http.StatusTooManyRequests) // 429

			_ = json.NewEncoder(w).Encode(RateLimitResponse{
				Error:      "rate_limit_exceeded",
				RetryAfter: retryAfter,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
