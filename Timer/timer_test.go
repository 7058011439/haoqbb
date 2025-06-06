package Timer

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAddOnceTimer(t *testing.T) {
	cost := NewTiming(Millisecond)
	var times50 int32
	var times1000 int32
	var times5000 int32

	count := 0
	AddRepeatTimer(50, func(timerId TimerID, args ...interface{}) {
		if count < 250 {
			AddRepeatTimer(1000, func(timerId TimerID, args ...interface{}) {
				atomic.AddInt32(&times1000, 1)
			})
			count++
		}
		atomic.AddInt32(&times50, 1)
	})

	AddRepeatTimer(5000, func(timerId TimerID, args ...interface{}) {
		atomic.AddInt32(&times5000, 1)
		Log.Debug("协程 %v 耗时 %v, 分别执行 %v %v %v", timerId, cost, atomic.LoadInt32(&times50), atomic.LoadInt32(&times1000), atomic.LoadInt32(&times5000))
		atomic.StoreInt32(&times50, 0)
		atomic.StoreInt32(&times1000, 0)
		atomic.StoreInt32(&times5000, 0)
		cost.ReStart()
	})
	select {}
}

func TestXXX(t *testing.T) {
	tick := time.NewTicker(time.Millisecond)
	lastTime := time.Now()
	for {
		select {
		case <-tick.C:
			fmt.Println(time.Now().Sub(lastTime).Microseconds())
			lastTime = time.Now()
		}
	}
}

func TestYYY(t *testing.T) {
	tick := time.NewTicker(time.Millisecond * 1)
	lastTime := time.Now()
	times := 0
	for {
		select {
		case <-tick.C:
			times++
			duration := time.Now().Sub(lastTime).Seconds()
			if duration > 30 {
				Log.Debug("耗时 %v 毫秒, 总计执行 %v 次", duration, times)
				lastTime = time.Now()
				times = 0
			}
		}
	}
}

func TestAddOnceTimer2(t *testing.T) {
	var endTimeList = map[TimerID]int64{}
	var durationTimeList = map[TimerID]int64{}
	var diffList = map[int64]int{}
	var mutex sync.RWMutex
	var group sync.WaitGroup
	totalTimers := 300000
	repeatTimes := 100000
	group.Add(totalTimers - repeatTimes)
	for i := 0; i < totalTimers; i++ {
		duration := rand.Intn(100000) + 10000
		if i < repeatTimes {
			AddRepeatTimer(int64(duration), func(timerId TimerID, args ...interface{}) {})
		} else {
			id := AddOnceTimer(int64(duration), func(timerId TimerID, args ...interface{}) {
				mutex.Lock()
				diff := time.Now().UnixMilli() - endTimeList[timerId]
				if timerId%1000 == 0 { // 这句话没鸟用, 单纯看个程序是否在运行
					Log.Debug("特殊定时器被执行:id = %v, 误差 = %v, 误差率 = %.2f%%", timerId, diff, float64(diff)/float64(durationTimeList[timerId])*100)
				}
				diffList[diff]++
				mutex.Unlock()
				group.Done()
			})
			if id%1000 == 0 { // 这句话没鸟用, 单纯看个程序是否在运行
				Log.Debug("添加特殊定时器:id = %v, 延时 = %v", id, duration)
			}
			mutex.Lock()
			endTimeList[id] = time.Now().UnixMilli() + int64(duration)
			durationTimeList[id] = int64(duration)
			mutex.Unlock()
		}
		// time.Sleep(time.Millisecond)
	}
	group.Wait()
	Log.Debug("误差分布如下:")
	for diff, times := range diffList {
		fmt.Println(diff, times)
	}
}
