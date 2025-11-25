// SPDX-FileCopyrightText: 2025 maxhash.io <dev@maxhash.io>
//
// SPDX-License-Identifier: AGPL-3.0-only

package maxhash

// StatsService defines the interface for retrieving pool and user statistics.
type StatsService interface {
	PoolStats() (PoolStats, error)
	UserStats(username string) (UserStats, error)
}

// PoolStats represents merged pool statistics.
type PoolStats struct {
	Runtime      int64   `json:"runtime"`
	LastUpdate   int64   `json:"lastupdate"`
	Users        int     `json:"users"`
	Workers      int     `json:"workers"`
	Idle         int     `json:"idle"`
	Disconnected int     `json:"disconnected"`
	Hashrate1m   string  `json:"hashrate1m"`
	Hashrate5m   string  `json:"hashrate5m"`
	Hashrate15m  string  `json:"hashrate15m"`
	Hashrate1hr  string  `json:"hashrate1hr"`
	Hashrate6hr  string  `json:"hashrate6hr"`
	Hashrate1d   string  `json:"hashrate1d"`
	Hashrate7d   string  `json:"hashrate7d"`
	Diff         float64 `json:"diff"`
	Accepted     int64   `json:"accepted"`
	Rejected     int64   `json:"rejected"`
	BestShare    int64   `json:"bestshare"`
	SPS1m        float64 `json:"sps1m"`
	SPS5m        float64 `json:"sps5m"`
	SPS15m       float64 `json:"sps15m"`
	SPS1h        float64 `json:"sps1h"`
}

// PoolStatusLine1 represents the first line of pool statistics.
type PoolStatusLine1 struct {
	Runtime      int64 `json:"runtime"`
	LastUpdate   int64 `json:"lastupdate"`
	Users        int   `json:"Users"`
	Workers      int   `json:"Workers"`
	Idle         int   `json:"Idle"`
	Disconnected int   `json:"Disconnected"`
}

// PoolStatusLine2 represents the second line of pool statistics.
type PoolStatusLine2 struct {
	Hashrate1m  string `json:"hashrate1m"`
	Hashrate5m  string `json:"hashrate5m"`
	Hashrate15m string `json:"hashrate15m"`
	Hashrate1hr string `json:"hashrate1hr"`
	Hashrate6hr string `json:"hashrate6hr"`
	Hashrate1d  string `json:"hashrate1d"`
	Hashrate7d  string `json:"hashrate7d"`
}

// PoolStatusLine3 represents the third line of pool statistics.
type PoolStatusLine3 struct {
	Diff      float64 `json:"diff"`
	Accepted  int64   `json:"accepted"`
	Rejected  int64   `json:"rejected"`
	BestShare int64   `json:"bestshare"`
	SPS1m     float64 `json:"SPS1m"`
	SPS5m     float64 `json:"SPS5m"`
	SPS15m    float64 `json:"SPS15m"`
	SPS1h     float64 `json:"SPS1h"`
}

// UserStats represents user statistics.
type UserStats struct {
	Hashrate1m  string   `json:"hashrate1m"`
	Hashrate5m  string   `json:"hashrate5m"`
	Hashrate1hr string   `json:"hashrate1hr"`
	Hashrate1d  string   `json:"hashrate1d"`
	Hashrate7d  string   `json:"hashrate7d"`
	LastShare   int64    `json:"lastshare"`
	Workers     int      `json:"workers"`
	Shares      int64    `json:"shares"`
	BestShare   float64  `json:"bestshare"`
	BestEver    int64    `json:"bestever"`
	Authorised  int64    `json:"authorised"`
	Worker      []Worker `json:"worker"`
}

// Worker represents a worker's statistics.
type Worker struct {
	WorkerName  string  `json:"workername"`
	Hashrate1m  string  `json:"hashrate1m"`
	Hashrate5m  string  `json:"hashrate5m"`
	Hashrate1hr string  `json:"hashrate1hr"`
	Hashrate1d  string  `json:"hashrate1d"`
	Hashrate7d  string  `json:"hashrate7d"`
	LastShare   int64   `json:"lastshare"`
	Shares      int64   `json:"shares"`
	BestShare   float64 `json:"bestshare"`
	BestEver    int64   `json:"bestever"`
}
