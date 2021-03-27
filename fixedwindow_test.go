package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
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
		requestInterval time.Duration
		expectedResults []*ratelimiter.Result
		expectedErrors  []error
	}{
		{
			name:            "requests within limit",
			limit:           5,
			duration:        time.Duration(5 * time.Second),
			requestInterval: time.Duration(100 * time.Millisecond),
			expectedResults: []*ratelimiter.Result{
				{
					Limit:     5,
					Remaining: 4,
				},
				{
					Limit:     5,
					Remaining: 3,
				},
				{
					Limit:     5,
					Remaining: 2,
				},
				{
					Limit:     5,
					Remaining: 1,
				},
				{
					Limit:     5,
					Remaining: 0,
				},
			},
			expectedErrors: []error{
				nil,
				nil,
				nil,
				nil,
				nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			mockClock := clock.NewMock()
			r := ratelimiter.NewFixedWindowRateLimiter(tc.limit, tc.duration, mockRepo, mockClock)

			for i := 0; i < len(tc.expectedResults); i++ {
				mockRepo.
					EXPECT().
					IncrementByKey(gomock.Any(), gomock.Eq("test_key"), gomock.Any()).
					Return(i+1, nil)

				res, err := r.Allow(context.Background(), "test_key")
				assert.Equal(t, tc.expectedResults[i], res)
				if tc.expectedErrors[i] == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.expectedErrors[i].Error())
				}
			}
		})
	}
}
