package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/viper"
	"github.com/yonasstephen/ratelimiter/examples/httpserver/server"
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

	httpServer := server.NewHTTPServer(server.Opts{
		Port:              port,
		RateLimitCount:    limit,
		RateLimitDuration: duration,
	})

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		oscall := <-c
		log.Println("system call:", oscall)
		cancel()
	}()

	if err := httpServer.Start(ctx); err != nil {
		log.Println("failed to start HTTP Server:", err)
	}
}
