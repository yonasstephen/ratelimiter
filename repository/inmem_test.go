package repository_test

import (
	"context"
	"sync"
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

// This test attempts to run multiple goroutines to increment by key
// which will modify the same key in the in-mem map. Without a mutex
// in the impelementation, it would raise the following warning.
//
// ==================
// WARNING: DATA RACE
// Read at 0x00c000094000 by goroutine 12:
//   github.com/yonasstephen/ratelimiter/repository.(*InMemRepository).IncrementByKey()
//       /Users/yonasstephen/go/src/github.com/yonasstephen/ratelimiter/repository/inmem.go:47 +0x215
//   github.com/yonasstephen/ratelimiter/repository_test.TestIncrementByKey_RaceCondition.func1()
//       /Users/yonasstephen/go/src/github.com/yonasstephen/ratelimiter/repository/inmem_test.go:65 +0xed

// Previous write at 0x00c000094000 by goroutine 9:
//   github.com/yonasstephen/ratelimiter/repository.(*InMemRepository).IncrementByKey()
//       /Users/yonasstephen/go/src/github.com/yonasstephen/ratelimiter/repository/inmem.go:49 +0xe5
//   github.com/yonasstephen/ratelimiter/repository_test.TestIncrementByKey_RaceCondition.func1()
//       /Users/yonasstephen/go/src/github.com/yonasstephen/ratelimiter/repository/inmem_test.go:65 +0xed

// Goroutine 12 (running) created at:
//   github.com/yonasstephen/ratelimiter/repository_test.TestIncrementByKey_RaceCondition()
//       /Users/yonasstephen/go/src/github.com/yonasstephen/ratelimiter/repository/inmem_test.go:63 +0x250
//   testing.tRunner()
//       /usr/local/Cellar/go/1.16.2/libexec/src/testing/testing.go:1194 +0x202

// Goroutine 9 (finished) created at:
//   github.com/yonasstephen/ratelimiter/repository_test.TestIncrementByKey_RaceCondition()
//       /Users/yonasstephen/go/src/github.com/yonasstephen/ratelimiter/repository/inmem_test.go:63 +0x250
//   testing.tRunner()
//       /usr/local/Cellar/go/1.16.2/libexec/src/testing/testing.go:1194 +0x202
// ==================
func TestIncrementByKey_RaceCondition(t *testing.T) {
	mockClock := clock.NewMock()
	ctx := context.Background()
	inMem := repository.NewInMemRepository()

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = inMem.IncrementByKey(ctx, "key1", mockClock.Now())
		}()
	}
	wg.Wait()
}
