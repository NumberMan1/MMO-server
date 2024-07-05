package timeunit

import (
	"time"
)

// TimeClock 记录地图的游戏时钟
type TimeClock struct {
	// StartTime 游戏开始的时间戳
	StartTime time.Time
	// Time 游戏的运行时间（秒）
	Time float64
	// DeltaTime 获取上一帧运行所用的时间（秒）
	DeltaTime float64
	// LastTick 记录最后一次tick的时间
	LastTick time.Time
}

func NewTimeClock() *TimeClock {
	return &TimeClock{
		StartTime: time.Now(),
	}
}

// Tick 由 summer/schedule调用，请不要自行调用，除非你知道自己在做什么！！！
func (tc *TimeClock) Tick() {
	now := time.Now()
	tc.Time = now.Sub(tc.StartTime).Seconds()
	if tc.LastTick.IsZero() {
		tc.LastTick = now
	}
	tc.DeltaTime = now.Sub(tc.LastTick).Seconds()
	tc.LastTick = now
}
