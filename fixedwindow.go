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

	// internal state
	hasExceeded bool
	curWindow   time.Time
}

// NewFixedWindowRateLimiter returns an instance of fixed window rate limiter.
// It takes limit & duration. The rate is defined as limit/duration. Example:
//
//   limit := 1
//   duration := 5*time.Second
//   // this gives us a rate of 1 request per 5 second window
//   rateLimiter := NewFixedWindowRateLimiter(limit, duration, repo, clock)
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
	if !r.curWindow.Equal(window) {
		r.curWindow = window
		r.hasExceeded = false
	}

	// if it has exceeded the limit before, do not increment store
	windowResetTime := window.Add(r.duration).Sub(now)
	if r.hasExceeded {
		return &Result{
			Limit:      r.limit,
			Remaining:  0,
			RetryAfter: windowResetTime,
			ResetAfter: windowResetTime,
		}, nil
	}

	// increment the request count in the store
	count, err := r.repo.IncrementByKey(ctx, key, window)
	if err != nil {
		return nil, errors.Wrap(err, "failed to increment repository")
	}

	// if exceeds the limit for the first time, flag hasExceeded
	if count > r.limit {
		r.hasExceeded = true
		return &Result{
			Limit:      r.limit,
			Remaining:  0,
			RetryAfter: windowResetTime,
			ResetAfter: windowResetTime,
		}, nil
	}

	return &Result{
		Limit:     r.limit,
		Remaining: r.limit - count,
	}, nil
}
