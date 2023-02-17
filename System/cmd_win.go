//go:build windows
// +build windows

package System

import (
	"syscall"
	"unsafe"
)

func SetTitle(title string) {
	kernel32, _ := syscall.LoadLibrary(`kernel32.dll`)
	sct, _ := syscall.GetProcAddress(kernel32, `SetConsoleTitleW`)
	strUtf16, _ := syscall.UTF16PtrFromString(title)
	syscall.Syscall(sct, 1, uintptr(unsafe.Pointer(strUtf16)), 0, 0)
	syscall.FreeLibrary(kernel32)
}
