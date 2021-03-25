package ratelimiter

import (
	"context"
	"time"
)

// RateLimiter is the interface of a rate limit module
type RateLimiter interface {
	// Allow increments the rate of the request for a given key and returns
	// information about the result of the request.
	Allow(ctx context.Context, key string) (*Result, error)
}

// Result embodies information about the current state of the rate limit
type Result struct {
	// Allowed indicates how many allowed request at time.Now
	Allowed int
	// Remaining indicates how many remaining request is allowed at time.Now
	Remaining int
	// RetryAfer indicates the duration that requester need to wait
	// until the request will be allowed. If the request is allowed,
	// a zero value will be returned
	RetryAfter time.Duration
}
