package Timer

import (
	"sync"
	"time"
)

type TimerID = int64
type Timestamp = int64

type TimerFun func(timerId TimerID, args ...interface{})

const (
	Repeat int = 0
	Once   int = 1
)

var (
	mutex                sync.Mutex
	newTimerId           = TimerID(100000)
	mapTimers            = make(map[Timestamp][]*timer) // 精确的时间戳触发队列
	mapTimerId           = make(map[TimerID]Timestamp)
	checkInterval        = int64(1) // 每隔10ms 检查一次
	lastExecuteTimeStamp int64      // 上一次执行的时间戳+1
)

type timer struct {
	timeId     TimerID
	expireTime Timestamp
	duration   Timestamp // 持续时间，用于 Repeat 模式
	callback   TimerFun
	args       []interface{}
	eType      int
}

func init() {
	setTimerResolution(1)
	go scheduler()
}

func now() int64 {
	return time.Now().UnixMilli()
}

func addTimer(t *timer) {
	t.expireTime = now() + t.duration
	mapTimers[t.expireTime] = append(mapTimers[t.expireTime], t)
	mapTimerId[t.timeId] = t.expireTime
}

// 高精度调度器
func scheduler() {
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Millisecond)
	lastExecuteTimeStamp = now()
	for {
		<-ticker.C
		executeDueTimers()
	}
}

func executeDueTimers() {
	mutex.Lock()
	defer mutex.Unlock()

	current := now()
	for lastExecuteTimeStamp <= current {
		for _, t := range mapTimers[lastExecuteTimeStamp] {
			go t.callback(t.timeId, t.args...)
			if t.eType == Repeat {
				t.expireTime = current + t.duration
				addTimer(t)
			} else {
				delete(mapTimerId, t.timeId)
			}
		}

		delete(mapTimers, lastExecuteTimeStamp)
		lastExecuteTimeStamp++
	}
}

func AddRepeatTimer(duration int64, callback TimerFun, args ...interface{}) TimerID {
	if duration < 1 {
		return -1
	}
	mutex.Lock()
	defer mutex.Unlock()

	newTimerId++
	t := &timer{
		timeId:   newTimerId,
		duration: duration,
		callback: callback,
		args:     args,
		eType:    Repeat,
	}
	addTimer(t)
	return newTimerId
}

func AddOnceTimer(duration int64, callback TimerFun, args ...interface{}) TimerID {
	if duration < 1 {
		return -1
	}
	mutex.Lock()
	defer mutex.Unlock()

	newTimerId++
	t := &timer{
		timeId:   newTimerId,
		duration: duration,
		callback: callback,
		args:     args,
		eType:    Once,
	}
	addTimer(t)
	return newTimerId
}

func CloseTimer(id TimerID) {
	mutex.Lock()
	defer mutex.Unlock()

	expireTs, ok := mapTimerId[id]
	if !ok {
		return
	}
	timers := mapTimers[expireTs]
	for i, t := range timers {
		if t.timeId == id {
			timers = append(timers[:i], timers[i+1:]...)
			break
		}
	}
	if len(timers) == 0 {
		delete(mapTimers, expireTs)
	} else {
		mapTimers[expireTs] = timers
	}
	delete(mapTimerId, id)
}

func Count() int {
	mutex.Lock()
	defer mutex.Unlock()
	return len(mapTimerId)
}
