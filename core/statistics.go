package core

import (
	"fmt"
	"sync/atomic"
	"time"
)

// GameStats holds current game statistics
var GameStats = Statistics{
	started: time.Now(),
}

// Statistics hold game stats
type Statistics struct {
	gamesPlayed      int64
	clientsConnected int64
	started          time.Time
}

// StatisticsInfo is serialized Statistics
type StatisticsInfo struct {
	GamesPlayed      int64
	ClientsConnected int64
	Uptime           string
}

// AddPlayedGame increments played games count
func (s *Statistics) AddPlayedGame() {
	atomic.AddInt64(&s.gamesPlayed, 1)
}

// OnClientConnected - call when client has connected
func (s *Statistics) OnClientConnected() {
	atomic.AddInt64(&s.clientsConnected, 1)
}

// OnClientDisconnected - call when client disconnected
func (s *Statistics) OnClientDisconnected() {
	atomic.AddInt64(&s.clientsConnected, -1)
}

// ToInfo get StatisticsInfo intance for current Statistics
func (s *Statistics) ToInfo() *StatisticsInfo {
	i := StatisticsInfo{
		ClientsConnected: s.clientsConnected,
		GamesPlayed:      s.gamesPlayed,
	}
	d := time.Now().Sub(s.started)
	i.Uptime = fmt.Sprintf("%02d:%02d:%02d",
		int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
	return &i
}
