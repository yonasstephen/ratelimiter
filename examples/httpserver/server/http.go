package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/yonasstephen/ratelimiter"
	"github.com/yonasstephen/ratelimiter/examples/httpserver/middleware"
	"github.com/yonasstephen/ratelimiter/repository"
)

// HTTPServer is a simple http server with rate limiter
type HTTPServer struct {
	opts Opts
}

// Opts stores the configuration options for running HTTP server
type Opts struct {
	Port              int
	RateLimitCount    int
	RateLimitDuration time.Duration
}

// NewHTTPServer instantiates a new HTTPServer object with the
// given server configurations
func NewHTTPServer(opts Opts) *HTTPServer {
	return &HTTPServer{opts: opts}
}

// Start runs the HTTPServer. This is a blocking function.
// To stop the server, send a cancel signal to the context.
func (s *HTTPServer) Start(ctx context.Context) error {
	// init dependencies
	inMemRepo := repository.NewInMemRepository()
	clock := clock.New()
	fixedWindowLimiter := ratelimiter.NewFixedWindowRateLimiter(s.opts.RateLimitCount, s.opts.RateLimitDuration, inMemRepo, clock)
	rateLimitMiddleware := middleware.NewRateLimiterMiddleware(fixedWindowLimiter, clock)

	// setup http handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handlePing)

	testHandler := http.HandlerFunc(handleTest)
	mux.Handle("/test", rateLimitMiddleware.AttachRateLimitMiddleware(testHandler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.opts.Port),
		Handler: mux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to listen and serve:", err)
		}
	}()

	log.Println("HTTP Server started on port", s.opts.Port)

	<-ctx.Done()

	// handle shutdown
	log.Println("Shutting down HTTP Server...", s.opts.Port)

	// set timeout for shutfown operation
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := srv.Shutdown(ctxShutDown)
	if err != nil {
		log.Fatal("server shutdown failed:", err)
	}

	log.Printf("server exited gracefully")
	return err
}

// handlePing is a health check endpoint
func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "request is successful!")
}
