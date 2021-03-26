package ratelimiter

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/pkg/errors"
	"github.com/yonasstephen/ratelimiter/repository"
)

// FixedWindowRateLimiter is an implementation of RateLimiter interface
// with a fixed window algorithm
type FixedWindowRateLimiter struct {
	clock    clock.Clock
	duration time.Duration
	limit    int
	repo     repository.Repository
}

// NewFixedWindowRateLimiter returns an instance of fixed window rate limiter.
// It takes limit & duration. The rate is defined as limit/duration. Example:
//
//   limit := 1
//   duration := 5*time.Second
//   // this gives us a rate of 1 request per 5 second window
//
func NewFixedWindowRateLimiter(limit int, duration time.Duration, repo repository.Repository, clock clock.Clock) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clock:    clock,
		duration: duration,
		limit:    limit,
		repo:     repo,
	}
}

// Allow increments the request rate of the given key for the current
// time window and returns the result
func (r *FixedWindowRateLimiter) Allow(ctx context.Context, key string) (*Result, error) {
	now := r.clock.Now()
	window := now.Truncate(r.duration)
	count, err := r.repo.IncrementByKey(ctx, key, window)
	if err != nil {
		return nil, errors.Wrap(err, "failed to increment repository")
	}

	allowed := 1
	// TODO: fix remaining calculation
	remaining := 0
	retryAfter := time.Duration()

	// if the request is not allowed
	if count > r.limit {
		allowed = 0
		retryAfter = window.Add(r.duration).Sub(now)
	}

	res := &Result{
		Allowed:    allowed,
		Remaining:  remaining,
		RetryAfter: retryAfter,
	}
	return res, nil
}
