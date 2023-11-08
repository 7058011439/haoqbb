package Timer

import (
	"sync"
	"time"
)

type TimerID = int64
type TimeWheel = int64

type TimerFun func(timerId TimerID, args ...interface{})

const (
	Repeat int = 0
	Once   int = 1
)

var mutex sync.Mutex

var newTimerId = TimerID(100000)
var mapTimerWheel = make(map[TimeWheel][]*timer) // map[nextTick]list(Timer)
var mapTimerId = make(map[TimerID]TimeWheel)     // map[timeId]nextTick
var currTick = time.Now().UnixNano() / TimeWheel(time.Millisecond)

type timer struct {
	timeId   TimerID
	duration TimeWheel
	funcName TimerFun
	args     []interface{}
	eType    int
}

func init() {
	setTimerResolution(1)
	go doSomething()
}

func addTime(time *timer) {
	nextTick := currTick + time.duration
	mapTimerWheel[nextTick] = append(mapTimerWheel[nextTick], time)
	mapTimerId[time.timeId] = nextTick
}

func doSomething() {
	tick := time.NewTicker(time.Millisecond)
	for {
		<-tick.C
		mutex.Lock()
		if times, ok := mapTimerWheel[currTick]; ok {
			for i := 0; i < len(times); i++ {
				currTimer := times[i]
				go currTimer.funcName(currTimer.timeId, currTimer.args...)
				if currTimer.eType == Repeat {
					addTime(currTimer)
				} else {
					delete(mapTimerId, currTimer.timeId)
				}
			}
			delete(mapTimerWheel, currTick)
		}
		currTick++
		mutex.Unlock()
	}
}

// AddRepeatTimer
// 添加一次循环定时器
// duration-时间间隔(毫秒)
// funcName-回调函数
// args-回调参数
func AddRepeatTimer(duration TimeWheel, funcName TimerFun, args ...interface{}) TimerID {
	mutex.Lock()
	defer mutex.Unlock()
	if duration < 1 {
		return -1
	}
	newTimerId++
	addTime(&timer{timeId: newTimerId, duration: duration, funcName: funcName, args: args, eType: Repeat})
	return newTimerId
}

func AddOnceTimer(duration TimeWheel, funcName TimerFun, args ...interface{}) TimerID {
	mutex.Lock()
	defer mutex.Unlock()
	if duration < 1 {
		return -1
	}
	newTimerId++
	addTime(&timer{timeId: newTimerId, duration: duration, funcName: funcName, args: args, eType: Once})
	return newTimerId
}

func CloseTimer(id TimerID) {
	mutex.Lock()
	defer mutex.Unlock()

	nextTick := mapTimerId[id]
	if times, ok := mapTimerWheel[nextTick]; ok {
		for i := 0; i < len(times); i++ {
			if times[i].timeId == id {
				times = append(times[:i], times[i+1:]...)
				break
			}
		}
		mapTimerWheel[nextTick] = times
	}
	delete(mapTimerId, id)
}

func Count() int {
	return len(mapTimerId)
}
