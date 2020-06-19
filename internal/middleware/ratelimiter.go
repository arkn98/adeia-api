package middleware

import (
	"net/http"

	log "adeia-api/internal/utils/logger"
	"adeia-api/internal/utils/ratelimiter"
)

// RateLimiter is a middleware that wraps the internal/utils/ratelimiter pkg.
// It limits the requests-per-second, per IP, using the provided rate-limiter.
// Preferably, it should be added to the global middleware.
func RateLimiter(limiter ratelimiter.RateLimiter) Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.GetLimiter(r.RemoteAddr).Allow() {
				log.Debugf("limiting request rate for IP: %q", r.RemoteAddr)
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
