package repository

//go:generate mockgen -package=mocks -destination=mocks/repository.go github.com/yonasstephen/ratelimiter/repository Repository

import (
	"context"
	"time"
)

// Repository interfaces the interaction with the underlying
// store where the rate limit data is persisted
type Repository interface {
	IncrementByKey(ctx context.Context, key string, window time.Time) (int, error)
}
