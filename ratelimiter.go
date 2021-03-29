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
	// Allowed is the number of requests that are allowed at time.Now().
	// Zero value means that the request is not allowed i.e. has exceeded
	// the rate limit threshold
	Allowed int

	// Limit is the limit that was used to get this result
	Limit int

	// Remaining indicates how many remaining request is allowed at time.Now
	Remaining int

	// RetryAfer indicates the duration that requester need to wait
	// until the request will be allowed. If the request is allowed,
	// a zero value will be returned
	RetryAfter time.Duration

	// ResetAfter indicates the duration that the requester need to wait
	// until the time moves to the next rate limit window and hence resetting
	// the count. You can also think of this as the time when Limit == Remaining.
	ResetAfter time.Duration
}
