package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"github.com/yonasstephen/ratelimiter/repository"
)

func TestIncrementByKey(t *testing.T) {
	mockClock := clock.NewMock()
	ctx := context.Background()

	// increment key1
	inMem := repository.NewInMemRepository()
	count, err := inMem.IncrementByKey(ctx, "key1", mockClock.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// increment key1 again
	count, err = inMem.IncrementByKey(ctx, "key1", mockClock.Now())
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// increment key1 again
	count, err = inMem.IncrementByKey(ctx, "key1", mockClock.Now())
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	// increment key2
	count, err = inMem.IncrementByKey(ctx, "key2", mockClock.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// increment key1 again
	count, err = inMem.IncrementByKey(ctx, "key1", mockClock.Now())
	assert.NoError(t, err)
	assert.Equal(t, 4, count)

	// increment key1 with different window t0+5min, should reset windowObj
	count, err = inMem.IncrementByKey(ctx, "key1", mockClock.Now().Add(5*time.Minute))
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// increment key1 with the first time window t0, restarted from 0
	count, err = inMem.IncrementByKey(ctx, "key1", mockClock.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
