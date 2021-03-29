package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/spf13/viper"
	"github.com/yonasstephen/ratelimiter"
	"github.com/yonasstephen/ratelimiter/examples/httpserver/middleware"
	"github.com/yonasstephen/ratelimiter/repository"
)

func main() {
	// read from env var - env var takes precedence
	// over env from config file
	viper.AutomaticEnv()

	// read from config file
	viper.SetConfigFile(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed to read config from .env:", err)
	}

	// load config
	viper.SetDefault("PORT", 8080)
	port := viper.GetInt("PORT")
	limit := viper.GetInt("RATE_LIMIT_COUNT")
	duration := viper.GetDuration("RATE_LIMIT_DURATION")

	// init dependencies
	inMemRepo := repository.NewInMemRepository()
	clock := clock.New()
	fixedWindowLimiter := ratelimiter.NewFixedWindowRateLimiter(limit, duration, inMemRepo, clock)
	rateLimitMiddleware := middleware.NewRateLimiterMiddleware(fixedWindowLimiter)

	// setup http handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handlePing)

	testHandler := http.HandlerFunc(handleTest)
	mux.Handle("/test", rateLimitMiddleware.AttachRateLimitMiddleware(testHandler))

	// serve http server
	log.Println("Running httpserver on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

// handlePing is a health check endpoint
func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "request is successful!")
}
