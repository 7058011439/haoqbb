//go:build linux

package Panic

import (
	"os"
	"syscall"
)

// RedirectStderr to the file passed in
func RedirectStderr(fileName string) (err error) {
	File.CreateDir(fileName)
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	err = syscall.Dup3(int(logFile.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		return
	}
	return
}
