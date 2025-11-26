// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

package http

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/mine-at/maxhash.io"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed index.html
var indexPageHTML string

//go:embed user.html
var userPageHTML string

// Server wraps http.Server.
type Server struct {
	StatsSvc maxhash.StatsService

	engine  *gin.Engine
	httpSrv *http.Server
	limiter *rate.Limiter
	store   persistence.CacheStore
	proxy   *httputil.ReverseProxy
}

// NewServer constructs a new Server with the given StatsService and optional listener.
func NewServer(statsSvc maxhash.StatsService) (*Server, error) {
	svr := &Server{
		StatsSvc: statsSvc,
	}

	// Setup rate limiter if enabled.
	if viper.GetBool("http.rate_limiter.enabled") {
		svr.limiter = rate.NewLimiter(
			rate.Limit(viper.GetFloat64("http.rate_limiter.rps")),
			viper.GetInt("http.rate_limiter.burst"),
		)

		slog.Debug("Rate limiter enabled ⏱️")
	}

	// Setup cache store if enabled.
	if viper.GetBool("http.cache.enabled") {
		ttl := viper.GetDuration("http.cache.ttl")
		if ttl <= 0 {
			return nil, errors.New("http.cache.ttl must be greater than 0 when caching is enabled")
		}

		svr.store = persistence.NewInMemoryStore(ttl)

		slog.Debug("Cache store enabled ⚡", "ttl", ttl)
	}

	// Setup reverse proxy if enabled.
	if viper.GetBool("http.proxy.enabled") {
		targetHostURL := viper.GetString("http.proxy.target_host_url")
		if targetHostURL == "" {
			return nil, errors.New("http.proxy.target_host_url is not set")
		}

		hostURL, err := url.Parse(targetHostURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse target host URL: %w", err)
		}

		// TODO: Make transport settings configurable via viper.
		transport := &http.Transport{
			MaxIdleConns:          1000,
			MaxConnsPerHost:       100,
			IdleConnTimeout:       time.Minute,
			ExpectContinueTimeout: 0,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}

		svr.proxy = httputil.NewSingleHostReverseProxy(hostURL)
		svr.proxy.Transport = transport

		slog.Debug("Reverse proxy enabled ✈️", "target_host_url", targetHostURL)
	}

	return svr, nil
}

// ListenAndServe will listen and serve on the server address. Blocks until the server is stopped.
func (s *Server) ListenAndServe() error {
	gin.SetMode(gin.ReleaseMode)

	// Enable debug mode if log_level is set to debug.
	if strings.EqualFold(viper.GetString("log_level"), "debug") {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	// Apply rate limiter middleware if enabled.
	if s.limiter != nil {
		engine.Use(s.limitHandler)
	}

	// API group.
	api := engine.Group("/v1")
	{
		// Use site-wide cache if store is configured.
		// Note: expire param is ignored since store was created with a default TTL.
		if s.store != nil {
			api.Use(cache.SiteCache(s.store, 0))

			slog.Debug("Site-wide caching enabled for API responses 🗄️")
		}

		api.GET("/pool", s.poolStatusHandler)
		api.GET("/users/:username", s.userHandler)
	}

	// Serve static files.
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("strip static prefix: %w", err)
	}
	engine.StaticFS("/static", http.FS(staticFS))

	// Page group.
	pages := engine.Group("/")
	{
		pages.GET("/users/:username", s.userPageHandler)
		pages.GET("/", s.indexPageHandler)
	}

	s.engine = engine

	s.httpSrv = &http.Server{
		Addr:           viper.GetString("http.addr"),
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	return s.httpSrv.ListenAndServe()
}

// GracefulShutdown will gracefully shutdown the server.
func (s *Server) GracefulShutdown(ctx context.Context) error {
	if s.httpSrv != nil {
		return s.httpSrv.Shutdown(ctx) //nolint: wrapcheck
	}

	return nil
}

func (s *Server) limitHandler(c *gin.Context) {
	if s.limiter != nil && !s.limiter.Allow() {
		c.String(http.StatusTooManyRequests, "Rate limit exceeded. Slow down! 🐢")
		c.Abort()
		return
	}

	c.Next()
}

func (s *Server) indexPageHandler(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", []byte(indexPageHTML))
}

func (s *Server) poolStatusHandler(c *gin.Context) {
	// If proxy is enabled, forward the request to the remote node instead of serving locally.
	if s.proxy != nil {
		s.proxy.ServeHTTP(c.Writer, c.Request)
		return
	}

	// Get pool stats from the local StatsService.
	stats, err := s.StatsSvc.PoolStats()
	if err != nil {
		slog.Error("Error getting pool stats", "error", err)
		c.String(500, "failed to get pool stats")
		return
	}

	c.JSON(200, stats)
}

func (s *Server) userPageHandler(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", []byte(userPageHTML))
}

func (s *Server) userHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" || !IsValidBitcoinAddress(username) {
		c.String(400, "invalid user path or Bitcoin address")
		return
	}

	// If proxy is enabled, forward the request to the remote node instead of serving locally.
	if s.proxy != nil {
		s.proxy.ServeHTTP(c.Writer, c.Request)
		return
	}

	// Get user stats from the local StatsService.
	userStats, err := s.StatsSvc.UserStats(username)
	if err != nil {
		slog.Error("Error getting user stats", "error", err)
		c.String(500, "failed to get user stats")
		return
	}

	c.JSON(200, userStats)
}

// IsValidBitcoinAddress returns true if the address is a valid Bitcoin address.
func IsValidBitcoinAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, nil)
	return err == nil
}
