package it_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/yonasstephen/ratelimiter/examples/httpserver/server"
)

type httpServerTestSuite struct {
	suite.Suite
	httpServer     *server.HTTPServer
	httpServerOpts server.Opts
	httpCtxCancel  context.CancelFunc
}

func TestHttpServerTestSuite(t *testing.T) {
	suite.Run(t, &httpServerTestSuite{})
}

func (s *httpServerTestSuite) SetupSuite() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed to read config from .env:", err)
	}

	// load config
	port := viper.GetInt("PORT")
	limit := viper.GetInt("RATE_LIMIT_COUNT")
	duration := viper.GetDuration("RATE_LIMIT_DURATION")
	s.httpServerOpts = server.Opts{
		Port:              port,
		RateLimitCount:    limit,
		RateLimitDuration: duration,
	}
	s.httpServer = server.NewHTTPServer(s.httpServerOpts)

	ctx, cancel := context.WithCancel(context.Background())
	s.httpCtxCancel = cancel
	go s.httpServer.Start(ctx)
}

func (s *httpServerTestSuite) TearDownSuite() {
	s.httpCtxCancel()
}

func (s *httpServerTestSuite) Test_PingEndpoint() {
	url := fmt.Sprintf("http://localhost:%d/ping", s.httpServerOpts.Port)
	resp, err := http.Get(url)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)
	defer resp.Body.Close()
	s.Equal("pong", string(body))
}

func (s *httpServerTestSuite) Test_TestEndpoint() {
	// TODO: make test more deterministic by using a mock clock
	url := fmt.Sprintf("http://localhost:%d/test", s.httpServerOpts.Port)

	var resp *http.Response
	var err error
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	for i := 0; i < s.httpServerOpts.RateLimitCount; i++ {
		resp, err = http.Get(url)
		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		s.NoError(err)
		s.Equal("request is successful!", string(body))
		s.Equal(strconv.Itoa(s.httpServerOpts.RateLimitCount), resp.Header.Get("RateLimit-Limit"))
		s.Equal(strconv.Itoa(s.httpServerOpts.RateLimitCount-(i+1)), resp.Header.Get("RateLimit-Remaining"))
		s.Len(resp.Header.Values("RateLimit-Reset-After"), 2)
	}

	// should hit rate limit
	resp, err = http.Get(url)
	s.NoError(err)
	s.Equal(http.StatusTooManyRequests, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)
	s.Equal("request has exceeded rate limit", string(body))
	s.Equal(strconv.Itoa(s.httpServerOpts.RateLimitCount), resp.Header.Get("RateLimit-Limit"))
	s.Equal("0", resp.Header.Get("RateLimit-Remaining"))
	// TODO: implement stricter RateLimit-Retry & Reset-After assertions
	s.Len(resp.Header.Values("RateLimit-Retry-After"), 2)
	s.Len(resp.Header.Values("RateLimit-Reset-After"), 2)

	// wait for the next rate limit window
	time.Sleep(2 * time.Second)
	resp, err = http.Get(url)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	body, err = ioutil.ReadAll(resp.Body)
	s.NoError(err)
	s.Equal("request is successful!", string(body))
	s.Equal(strconv.Itoa(s.httpServerOpts.RateLimitCount), resp.Header.Get("RateLimit-Limit"))
	s.Equal(strconv.Itoa(s.httpServerOpts.RateLimitCount-1), resp.Header.Get("RateLimit-Remaining"))
	s.Len(resp.Header.Values("RateLimit-Reset-After"), 2)
}
