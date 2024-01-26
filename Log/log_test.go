package Log

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	times := 100000
	runTimes := int32(0)
	go func() {
		for i := 0; i < times; i++ {
			ErrorLog("Error %v", i)
			atomic.AddInt32(&runTimes, 1)
		}
	}()

	go func() {
		for i := 0; i < times; i++ {
			WarningLog("Warning %v", i)
			atomic.AddInt32(&runTimes, 1)
		}
	}()

	go func() {
		for i := 0; i < times; i++ {
			Log("Log %v", i)
			atomic.AddInt32(&runTimes, 1)
		}
	}()

	go func() {
		for i := 0; i < times; i++ {
			Debug("Debug %v", i)
			atomic.AddInt32(&runTimes, 1)
		}
	}()

	go func() {
		for i := 0; i < times; i++ {
			fmt.Println("Println ", i)
			atomic.AddInt32(&runTimes, 1)
		}
	}()

	for {
		if runTimes == int32(times*5) {
			break
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	//time.Sleep(time.Second)
	//
	//for i := 0; i < 100000; i++ {
	//	ErrorLog("Error %v", i)
	//}
	//
	//for i := 0; i < 100000; i++ {
	//	WarningLog("Warning %v", i)
	//}
	//
	//for i := 0; i < 100000; i++ {
	//	Log("Log %v", i)
	//}
	//
	//for i := 0; i < 100000; i++ {
	//	Debug("Debug %v", i)
	//}

	//time.Sleep(time.Hour)
}
