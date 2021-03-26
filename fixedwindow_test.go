package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/yonasstephen/ratelimiter"
	mock "github.com/yonasstephen/ratelimiter/repository/mock"
)

func TestAllow(t *testing.T) {
	testCases := []struct {
		name            string
		limit           int
		duration        time.Duration
		numOfRequests   int
		requestInterval time.Duration
		expectedResults []*ratelimiter.Result
		expectedErrors  []error
	}{
		{
			name:            "requests within limit",
			limit:           1,
			duration:        time.Duration(5 * time.Second),
			numOfRequests:   4,
			requestInterval: time.Duration(100 * time.Millisecond),
			expectedResults: []*ratelimiter.Result{
				&ratelimiter.Result{
					Allowed:   1,
					Remaining: 4,
				},
			},
			expectedErrors: []error{
				nil,
				nil,
				nil,
				nil,
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
				assert.Equal(t, tc.expectedResults[i], res)
				assert.EqualError(t, err, tc.expectedErrors[i].Error())
			}
		})
	}
}
