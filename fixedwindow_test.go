package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/yonasstephen/ratelimiter"
	mock "github.com/yonasstephen/ratelimiter/mock"
)

func TestAllow(t *testing.T) {
	testCases := []struct {
		name            string
		limit           int
		duration        time.Duration
		numOfRequests   int
		requestInterval time.Duration
		expectedResult  []*ratelimiter.Result
		expectedError   error
	}{
		{
			name:            "requests within limit",
			limit:           1,
			duration:        time.Duration(5 * time.Second),
			numOfRequests:   4,
			requestInterval: time.Duration(100 * time.Millisecond),
			expectedResult: []*ratelimiter.Result{
				&ratelimiter.Result{
					Allowed:   1,
					Remaining: 4,
				},
			},
			// TODO: add expected errors as slice
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			r := ratelimiter.NewFixedWindowRateLimiter(tc.limit, tc.duration, mockRepo)
			// TODO: mock IncrementByKey()

			for i := 0; i < tc.numOfRequests; i++ {
				res, err := r.Allow(context.Background(), "test_key")
				assert.Equal(t, tc.expectedResult[i], res)
				assert.ErrorEqual(t, err, tc.ExpectedError)
			}
		})
	}
}
