package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/yonasstephen/ratelimiter"
)

// RateLimitMiddleware is a http middleware for applying rate limit to an API.
// If the limit is exceeded, a status code 429 is returned. If an error is
// encountered while checking the limit, a status code 500 is returned. It also
// set the headers with rate limit information such as limit, retry after, and
// reset after. Ref: https://tools.ietf.org/id/draft-polli-ratelimit-headers-00.html
type RateLimitMiddleware struct {
	limiter ratelimiter.RateLimiter
}

// NewRateLimiterMiddleware instantiates a new rate limiter middleware
func NewRateLimiterMiddleware(limiter ratelimiter.RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
	}
}

// AttachRateLimitMiddleware wraps the passed http handler with a rate limiter middleware.
// The passed handler is only called if the rate limit threshold has not exceeded yet.
func (m *RateLimitMiddleware) AttachRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := m.limiter.Allow(r.Context(), "test")
		if err != nil {
			log.Println("failed to check rate limit:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to check rate limit: %s", err.Error())))
			return
		}

		if res.Allowed == 0 {
			// request is not allowed
			w.Header().Set("RateLimit-Limit", strconv.Itoa(res.Limit))
			w.Header().Set("RateLimit-Retry-After", res.RetryAfter.String())
			w.Header().Set("RateLimit-Reset-After", res.ResetAfter.String())
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("request has exceeded rate limit"))
			return
		}

		// request is allowed
		next.ServeHTTP(w, r)
	})
}
