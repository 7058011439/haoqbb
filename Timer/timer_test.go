package Timer

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
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
		Log.Debug("协程 %v 耗时 %v, 总计执行 %v 次, 分别执行 %v %v %v", timerId, cost, GetRunTimes(), atomic.LoadInt32(&times50), atomic.LoadInt32(&times1000), atomic.LoadInt32(&times5000))
		atomic.StoreInt32(&times50, 0)
		atomic.StoreInt32(&times1000, 0)
		atomic.StoreInt32(&times5000, 0)
		ResetRunTimes()
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
