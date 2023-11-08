// set_timer_windows.go
//go:build windows
// +build windows

package Timer

import "syscall"

var (
	winmm           = syscall.NewLazyDLL("winmm.dll")
	timeBeginPeriod = winmm.NewProc("timeBeginPeriod")
	timeEndPeriod   = winmm.NewProc("timeEndPeriod")
)

func setTimerResolution(resolution int) error {
	ret, _, _ := timeBeginPeriod.Call(uintptr(resolution))
	if ret != 0 {
		return syscall.Errno(ret)
	}
	return nil
}

func restoreTimerResolution() error {
	ret, _, _ := timeEndPeriod.Call(uintptr(1)) // 1 表示还原分辨率
	if ret != 0 {
		return syscall.Errno(ret)
	}
	return nil
}
