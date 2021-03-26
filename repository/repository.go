package repository

import (
	"context"
	"time"
)

// Repository interfaces the interaction with the underlying
// store where the rate limit data is persisted
type Repository interface {
	IncrementByKey(ctx context.Context, key string, window time.Time) (int, error)
}
