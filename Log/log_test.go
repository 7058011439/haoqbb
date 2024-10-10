package Log

import (
	"sync"
	"testing"
)

func TestPrint(t *testing.T) {
	Debug("debug info")
	Log("log info")
	Warn("warn info")
	Error("error info")
	Fatal("fatal info")
}

func TestAll(t *testing.T) {
	times := 1000000
	var waitGroup sync.WaitGroup
	SetPrintLevel(LevelFatal)
	waitGroup.Add(1)
	go func() {
		for i := 0; i < times; i++ {
			Error("Error %v", i)
		}
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		for i := 0; i < times; i++ {
			Warn("Warn %v", i)
		}
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		for i := 0; i < times; i++ {
			Log("Log %v", i)
		}
		waitGroup.Done()
	}()

	/*
		waitGroup.Add(1)
		go func() {
			for i := 0; i < times; i++ {
				Debug("Debug %v", i)
			}
			waitGroup.Done()
		}()
	*/
	waitGroup.Wait()

	Fatal("Fatal msg")
}
