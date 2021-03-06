[![Go Reference](https://pkg.go.dev/badge/github.com/yonasstephen/ratelimiter.svg)](https://pkg.go.dev/github.com/yonasstephen/ratelimiter)
![Build Status](https://img.shields.io/github/workflow/status/yonasstephen/ratelimiter/Go)
[![Coverage Status](https://coveralls.io/repos/github/yonasstephen/ratelimiter/badge.svg?branch=master)](https://coveralls.io/github/yonasstephen/ratelimiter?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/yonasstephen/ratelimiter)](https://goreportcard.com/report/github.com/yonasstephen/ratelimiter)

<img src="./gopher.png" height="180" alt="gopher logo credit: https://quasilyte.dev/gopherkon">

# Rate Limiter
This package provides extensible rate limiter module in Go. There are 2 main extensible points:
1. Rate limit algorithm (fixed window, sliding window, leaky bucket, etc.)
2. Data store - to track request count by key (in-memory, Redis, Hazelcast, etc.)

## Supported Algorithm
### Fixed Window
This is the simplest algorithm for rate limiting. It divides the time into fixed window. For example a rate of 5 per 12 minutes gives us the following time windows:
1. hh:00 - hh:11
2. hh:12 - hh:23
3. hh:24 - hh:35
4. hh:36 - hh:47
5. hh:48 - hh:59

Where hh is any hours in the clock. This algorithm is susceptible to spike near the window boundaries. For instance 5 requests at hh:11 and 5 requests at hh:12 are allowed because they happen to fall on 2 windows although if you see it without the windows, you are allowing 10 requests within 2 minutes.

## Supported Data Store
### In-memory
This is the simplest storage i.e. relying on in-mem data structure that is map to keep track of the request count. This is susceptible to data loss when the app restarts because the data is not persisted on disk.

## How to use
```
go get github.com/yonasstephen/ratelimiter
```
Use it in your code
```go
import github.com/yonasstephen/ratelimiter

func main() {
    repo := mocks.NewInMemRepository()
    clock := clock.Clock()
    // setting the limit to 5 per minute
    r := ratelimiter.NewFixedWindowRateLimiter("5", time.Minute, repo, clock)

    // increasing the rate limit count for user_123
    res, err := r.Allow(context.Background(), "user_123")
    if err != nil {
        fmt.Fatal("failed to check rate limit")
    }
    fmt.Println(res)
}
```
There exists an example on how to use the ratelimiter module as a HTTP middleware as well in the [examples/httpserver](https://github.com/yonasstephen/ratelimiter/tree/master/examples/httpserver) folder.

## What's next
These are future improvements that can be made on this module:
- [x] Thread-safe implementation of in-mem repository
- [x] Integration test for examples/httpserver
- [ ] Sliding window algorithm implementation
- [ ] Redis repository implementation

## Contributing
Run tests
```
make test
```
If you make any changes to interface contract, you can run go generate to regenerate the mocks
```
make generate
```
