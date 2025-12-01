// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	stdlib_os "os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mine-at/maxhash.io/http"
	"github.com/mine-at/maxhash.io/os"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Exit code to return on exit.
var code int

// Path to config file.
var configPath string

func main() {
	defer stdlib_os.Exit(code)

	pflag.StringVar(&configPath, "config", "config.yaml", "Config file path")
	pflag.Parse()

	if err := initConfig(); err != nil {
		slog.Error("Error initializing config", "err", err)
		code = 1
		return
	}

	setupLogger()

	statsSvc, err := os.NewStatsService()
	if err != nil {
		slog.Error("Error creating stats service", "err", err)
		code = 1
		return
	}

	svr, err := http.NewServer(statsSvc)
	if err != nil {
		slog.Error("Error creating server", "err", err)
		code = 1
		return
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return svr.ListenAndServe()
	})

	eg.Go(func() error {
		ctx, cancel := signal.NotifyContext(ctx, stdlib_os.Interrupt, syscall.SIGTERM)
		defer cancel()

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		return svr.GracefulShutdown(shutdownCtx) //nolint: noctx, contextcheck
	})

	slog.Info("HTTP server listening", "addr", viper.GetString("http.addr"))

	if err := eg.Wait(); err != nil && !strings.Contains(err.Error(), "received signal") {
		slog.Error("HTTP server error", "err", err)
		code = 1
	}

	slog.Info("HTTP server shutdown")
}

func initConfig() error {
	// Use config file from the flag if set.
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Set default config values.
	viper.SetDefault("log_level", "info")
	viper.SetDefault("http.addr", "[::]:8080")
	viper.SetDefault("ckpool.log_dir", "/var/log/ckpool")
	viper.SetDefault("http.rate_limiter.enabled", false)
	viper.SetDefault("http.rate_limiter.rps", 5.0)
	viper.SetDefault("http.rate_limiter.burst", 10)
	viper.SetDefault("http.cache.enabled", true)
	viper.SetDefault("http.cache.ttl", time.Minute)
	viper.SetDefault("ckpool.http.proxy.enabled", false)
	viper.SetDefault("ckpool.http.proxy.target_host_url", "http://main.maxhash.io:8080")

	// Read in environment variables that match.
	viper.AutomaticEnv()

	// Tries to read the config file, if not found creates one with defaults.
	var fileLookupError viper.ConfigFileNotFoundError
	if err := viper.ReadInConfig(); err != nil {
		if stdlib_os.IsNotExist(err) || errors.As(err, &fileLookupError) {
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				return fmt.Errorf("write default config to %q: %w", configPath, err)
			}
			slog.Warn("Wrote default config", "path", configPath)
		} else {
			return fmt.Errorf("failed to read config from %q: %w", configPath, err)
		}
	}

	return nil
}

func setupLogger() {
	logLevel := slog.LevelInfo

	lvlStr := viper.GetString("log_level")

	if lvlStr != "" {
		switch lvlStr {
		case "DEBUG", "debug":
			logLevel = slog.LevelDebug
		case "INFO", "info":
			logLevel = slog.LevelInfo
		case "WARNING", "warning":
			logLevel = slog.LevelWarn
		case "ERROR", "error":
			logLevel = slog.LevelError
		}
	}

	h := slog.NewTextHandler(stdlib_os.Stdout, &slog.HandlerOptions{Level: logLevel})

	slog.SetDefault(slog.New(h))

	slog.Debug("Logger initialized", "level", logLevel.String())
}
