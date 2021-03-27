package ratelimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yonasstephen/ratelimiter"
	"github.com/yonasstephen/ratelimiter/repository/mocks"
)

type repoExpectation struct {
	timeWindow string // time string in RFC3339
	count      int
	err        error
}

type timeMatcher struct {
	expectedTime time.Time
}

func matchesTime(tm time.Time) gomock.Matcher {
	return &timeMatcher{expectedTime: tm}
}

// Matches returns whether x is a match.
func (m *timeMatcher) Matches(x interface{}) bool {
	v, ok := x.(time.Time)
	if !ok {
		return false
	}
	return m.expectedTime.Equal(v)
}

// String describes what the matcher matches.
func (m *timeMatcher) String() string {
	return "matches time using time.Equal() where different timezone can match if the datetime are equal when converted to the same timezone"
}

func TestAllow(t *testing.T) {
	testCases := []struct {
		name               string
		limit              int
		duration           time.Duration
		requestInterval    time.Duration
		expectedRepoReturn []repoExpectation
		expectedResults    []*ratelimiter.Result
		expectedErrors     []error
	}{
		{
			name:            "requests within limit",
			limit:           5,
			duration:        time.Duration(5 * time.Second),
			requestInterval: time.Duration(100 * time.Millisecond),
			expectedRepoReturn: []repoExpectation{
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      1,
					err:        nil,
				},
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      2,
					err:        nil,
				},
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      3,
					err:        nil,
				},
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      4,
					err:        nil,
				},
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      5,
					err:        nil,
				},
			},
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
		{
			name:            "requests exceeds limit",
			limit:           2,
			duration:        time.Duration(4 * time.Second),
			requestInterval: time.Second,
			expectedRepoReturn: []repoExpectation{
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      1,
					err:        nil,
				},
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      2,
					err:        nil,
				},
				{
					timeWindow: "1970-01-01T00:00:00Z",
					count:      3,
					err:        nil,
				},
				// should not make any repo call because hasExceed=true
				{},
				// should make repo call because window has changed
				{
					timeWindow: "1970-01-01T00:00:04Z",
					count:      1,
					err:        nil,
				},
			},
			expectedResults: []*ratelimiter.Result{
				{
					Limit:     2,
					Remaining: 1,
				},
				{
					Limit:     2,
					Remaining: 0,
				},
				// has exceeded limit
				{
					Limit:      2,
					Remaining:  0,
					RetryAfter: time.Duration(2 * time.Second),
					ResetAfter: time.Duration(2 * time.Second),
				},
				{
					Limit:      2,
					Remaining:  0,
					RetryAfter: time.Duration(1 * time.Second),
					ResetAfter: time.Duration(1 * time.Second),
				},
				// should reset on next window
				{
					Limit:     2,
					Remaining: 1,
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

			mockRepo := mocks.NewMockRepository(ctrl)
			mockClock := clock.NewMock()
			r := ratelimiter.NewFixedWindowRateLimiter(tc.limit, tc.duration, mockRepo, mockClock)

			for i := 0; i < len(tc.expectedResults); i++ {
				// expect repo call only if repoExpectation is not empty
				emptyRepoExpectation := repoExpectation{}
				if tc.expectedRepoReturn[i] != emptyRepoExpectation {
					expectedWindow, err := time.Parse(time.RFC3339, tc.expectedRepoReturn[i].timeWindow)
					require.NoError(t, err)

					mockRepo.
						EXPECT().
						IncrementByKey(gomock.Any(), gomock.Eq("test_key"), matchesTime(expectedWindow)).
						Return(tc.expectedRepoReturn[i].count, tc.expectedRepoReturn[i].err)
				}

				res, err := r.Allow(context.Background(), "test_key")
				assert.Equal(t, tc.expectedResults[i], res)
				if tc.expectedErrors[i] == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.expectedErrors[i].Error())
				}

				// move forward the time by specified interval
				mockClock.Add(tc.requestInterval)
			}
		})
	}
}
