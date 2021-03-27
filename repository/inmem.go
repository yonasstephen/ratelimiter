package repository

import (
	"context"
	"time"
)

// InMemRepository is repository implementation with a Go in-mem map.
// Note that on server restarts, the rate limit will be reset due to
// in-mem approach.
type InMemRepository struct {
	store map[string]windowObj
}

type windowObj struct {
	time  time.Time
	count int
}

// NewInMemRepository returns a new instance of in-mem repository
func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		store: map[string]windowObj{},
	}
}

// IncrementByKey increases the request count for the given key and
// current window by 1. If it exceeds the limit...
func (r *InMemRepository) IncrementByKey(ctx context.Context, key string, window time.Time) (int, error) {
	// TODO: implement mutex lock to support multi thread
	w, ok := r.store[key]
	if !ok || !w.time.Equal(window) {
		r.store[key] = windowObj{
			time:  window,
			count: 1,
		}
		return 1, nil
	}

	w.count++
	return w.count, nil
}