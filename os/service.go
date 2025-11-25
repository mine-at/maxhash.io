// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

package os

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mine-at/maxhash.io"
	"github.com/spf13/viper"
)

// StatsService implements maxhash.StatsService using the filesystem.
type StatsService struct {
	LogDir string
}

// NewStatsService constructs a new StatsService using viper config.
func NewStatsService() (*StatsService, error) {
	logDir := viper.GetString("ckpool.log_dir")
	if logDir == "" {
		return nil, errors.New("ckpool.log_dir is not set")
	}

	return &StatsService{
		LogDir: logDir,
	}, nil
}

var _ maxhash.StatsService = (*StatsService)(nil)

func (s *StatsService) PoolStats() (maxhash.PoolStats, error) {
	path := filepath.Join(s.LogDir, "pool", "pool.status")

	bytes, err := os.ReadFile(path)
	if err != nil {
		return maxhash.PoolStats{}, fmt.Errorf("read pool.status file: %w", err)
	}

	lines := strings.Split(string(bytes), "\n")
	if len(lines) < 3 {
		return maxhash.PoolStats{}, errors.New("invalid pool.status file format")
	}

	var (
		l1 maxhash.PoolStatusLine1
		l2 maxhash.PoolStatusLine2
		l3 maxhash.PoolStatusLine3
	)
	if err := json.Unmarshal([]byte(lines[0]), &l1); err != nil {
		return maxhash.PoolStats{}, fmt.Errorf("failed to unmarshal pool.status line 1: %w", err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &l2); err != nil {
		return maxhash.PoolStats{}, fmt.Errorf("failed to unmarshal pool.status line 2: %w", err)
	}
	if err := json.Unmarshal([]byte(lines[2]), &l3); err != nil {
		return maxhash.PoolStats{}, fmt.Errorf("failed to unmarshal pool.status line 3: %w", err)
	}

	return maxhash.PoolStats{
		LastUpdate:   l1.LastUpdate,
		Users:        l1.Users,
		Workers:      l1.Workers,
		Idle:         l1.Idle,
		Disconnected: l1.Disconnected,
		Hashrate1m:   l2.Hashrate1m,
		Hashrate5m:   l2.Hashrate5m,
		Hashrate15m:  l2.Hashrate15m,
		Hashrate1hr:  l2.Hashrate1hr,
		Hashrate6hr:  l2.Hashrate6hr,
		Hashrate1d:   l2.Hashrate1d,
		Hashrate7d:   l2.Hashrate7d,
		Diff:         l3.Diff,
		Accepted:     l3.Accepted,
		Rejected:     l3.Rejected,
		BestShare:    l3.BestShare,
		SPS1m:        l3.SPS1m,
		SPS5m:        l3.SPS5m,
		SPS15m:       l3.SPS15m,
		SPS1h:        l3.SPS1h,
	}, nil
}

func (s *StatsService) UserStats(username string) (maxhash.UserStats, error) {
	path := filepath.Join(s.LogDir, "users", username)

	bytes, err := os.ReadFile(path)
	if err != nil {
		return maxhash.UserStats{}, fmt.Errorf("read user stats file: %w", err)
	}

	var stats maxhash.UserStats
	if err := json.Unmarshal(bytes, &stats); err != nil {
		return maxhash.UserStats{}, fmt.Errorf("unmarshal user stats file: %w", err)
	}

	return stats, nil
}
