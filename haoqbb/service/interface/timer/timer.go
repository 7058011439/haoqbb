package ITimer

import (
	"github.com/7058011439/haoqbb/Timer"
	"runtime"
)

type ITimer interface {
	SetRepeatTimer(duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID
	SetOnceTimer(duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID
	SetOffTimer(id Timer.TimerID)
	GetTimerCount() int
	GetName() string
}

type serviceTimer struct {
	i map[string]ITimer
}

var timer = serviceTimer{i: make(map[string]ITimer)}

func SetTimerAgent(d ITimer) {
	timer.i[d.GetName()] = d
}

func SetRepeatTimer(serviceName string, duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID {
	if runtime.GOOS == "windows" {
		duration /= 1
		if duration < 1 {
			duration = 1
		}
	}
	return timer.i[serviceName].SetRepeatTimer(duration, funcName, args...)
}

func SetOnceTimer(serviceName string, duration Timer.TimeWheel, funcName Timer.TimerFun, args ...interface{}) Timer.TimerID {
	if runtime.GOOS == "windows" {
		duration /= 1
		if duration < 1 {
			duration = 1
		}
	}
	return timer.i[serviceName].SetOnceTimer(duration, funcName, args...)
}

func SetOffTimer(serviceName string, id Timer.TimerID) {
	timer.i[serviceName].SetOffTimer(id)
}

func GetTimerCount(serviceName string) int {
	return timer.i[serviceName].GetTimerCount()
}
