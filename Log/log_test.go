package Log

import (
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	go func() {
		for i := 0; i < 100000; i++ {
			ErrorLog("Error %v", i)
		}
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			WarningLog("Warning %v", i)
		}
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			Log("Log %v", i)
		}
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			Debug("Debug %v", i)
		}
	}()

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

	time.Sleep(time.Hour)
}
