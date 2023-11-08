// set_timer_linux.go
//go:build linux
// +build linux

package Timer

func setTimerResolution(resolution int) error {
	// do nothing
	return nil
}

func restoreTimerResolution() error {
	// do nothing
	return nil
}
