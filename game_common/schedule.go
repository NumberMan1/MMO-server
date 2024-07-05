package game_common

import (
	"context"
	"github.com/NumberMan1/MMO-server/game_common/timeunit"
	"github.com/NumberMan1/common/ns"
	"github.com/NumberMan1/common/ns/singleton"
	"reflect"
	"sync"
	"time"
)

var (
	singleSchedule = singleton.Singleton{}
)

// GetScheduleInstance 获取时钟调度器单例
func GetScheduleInstance() *Schedule {
	instance, _ := singleton.GetOrDo[*Schedule](&singleSchedule, func() (*Schedule, error) {
		return NewSchedule(), nil
	})
	return instance
}

type task struct {
	TaskMethod   func() error
	StartTime    time.Time
	Interval     time.Duration
	RepeatCount  int
	currentCount int
	lastTick     time.Time //上一次执行开始的时间
	Completed    bool      //是否已经执行完毕
}

func (t *task) Run() {
	t.lastTick = time.Now()
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	t.TaskMethod()
	t.currentCount += 1
	if (t.currentCount == t.RepeatCount) && (t.RepeatCount != 0) {
		t.Completed = true
	}
}

func (t *task) ShouldRun() bool {
	if (t.currentCount == t.RepeatCount) && (t.RepeatCount != 0) {
		return false
	}
	now := time.Now()
	if (now.After(t.StartTime)) && ((now.Sub(t.lastTick)) >= t.Interval) {
		return true
	}
	return false
}

// repeatCount 为0表示无限重复
func newTask(taskMethod func() error, startTime time.Time, interval time.Duration, repeatCount int) *task {
	return &task{
		TaskMethod:   taskMethod,
		StartTime:    startTime,
		Interval:     interval,
		RepeatCount:  repeatCount,
		currentCount: 0,
		lastTick:     time.Time{},
		Completed:    false,
	}
}

type Schedule struct {
	tasks       []*task
	addQueue    *ns.TSQueue[*task]
	removeQueue *ns.TSQueue[func() error]
	isStart     bool
	stop        chan struct{}
	fps         int // 每秒帧数
	ticker      *time.Ticker
	next        time.Time //下一帧执行的时间
	mutex       sync.Mutex
	timeClock   *timeunit.TimeClock
}

func NewSchedule() *Schedule {
	r := &Schedule{
		tasks:       make([]*task, 0),
		addQueue:    ns.NewTSQueue[*task](),
		removeQueue: ns.NewTSQueue[func() error](),
		isStart:     false,
		stop:        make(chan struct{}, 1),
		fps:         80,
		ticker:      nil,
		mutex:       sync.Mutex{},
		timeClock:   timeunit.NewTimeClock(),
	}
	return r
}

func (s *Schedule) Start(ctx context.Context) *Schedule {
	if !s.isStart {
		s.isStart = true
		s.ticker = time.NewTicker(45 * time.Millisecond)
		go s.execute()
	}
	return s
}

func (s *Schedule) Clock() *timeunit.TimeClock {
	return s.timeClock
}

func (s *Schedule) Stop() *Schedule {
	if s.isStart {
		s.ticker.Stop()
		s.isStart = false
		s.stop <- struct{}{}
	}
	return s
}

func (s *Schedule) AddTask(action func() error, timeValue time.Duration, repeatCount int) {
	startTime := time.Now()
	t := newTask(action, startTime, timeValue, repeatCount)
	s.addQueue.Push(t)
}

func (s *Schedule) RemoveTask(action func() error) {
	s.removeQueue.Push(action)
}

// Update 每帧都会执行
func (s *Schedule) Update(action func() error) {
	t := newTask(action, time.Now(), time.Millisecond, 0)
	s.addQueue.Push(t)
}

func (s *Schedule) execute() {
	for {
		select {
		case <-s.ticker.C:
			// tick间隔
			interval := time.Second / time.Duration(s.fps)
			startTime := time.Now()
			if startTime.Before(s.next) {
				continue
			}
			s.next = startTime.Add(interval)
			s.timeClock.Tick()
			s.mutex.Lock()
			//移除队列
			for item := s.removeQueue.Pop(); item != nil; item = s.removeQueue.Pop() {
				for i, task := range s.tasks {
					if reflect.ValueOf(task.TaskMethod) == reflect.ValueOf(item) {
						s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
						break
					}
				}
			}
			// 移除完毕的任务
			for i, task := range s.tasks {
				if task.Completed {
					s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
					break
				}
			}
			// 添加队列任务
			for item := s.addQueue.Pop(); item != nil; item = s.addQueue.Pop() {
				s.tasks = append(s.tasks, item)
			}
			// 执行任务
			for _, task := range s.tasks {
				if task.ShouldRun() {
					task.Run()
				}
			}
			s.mutex.Unlock()
		case <-s.stop:
			return
		}
	}
}
