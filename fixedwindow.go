package ratelimiter

import (
	"context"
	"time"

	"github.com/yonasstephen/ratelimiter/repository"
)

// FixedWindowRateLimiter is an implementation of RateLimiter interface
// with a fixed window algorithm
type FixedWindowRateLimiter struct {
	limit    int
	duration time.Duration
	repo     repository.Repository
}

// NewFixedWindowRateLimiter returns an instance of fixed window rate limiter
func NewFixedWindowRateLimiter(limit int, duration time.Duration, repo repository.Repository) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		limit:    limit,
		duration: duration,
		repo:     repo,
	}
}

// Allow increments the request rate of the given key for the current
// time window and returns the result
func (r *FixedWindowRateLimiter) Allow(ctx context.Context, key string) (*Result, error) {
	return nil, nil
}
