//go:build windows
// +build windows

package singleton

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutexW = kernel32.NewProc("CreateMutexW")
	procCloseHandle  = kernel32.NewProc("CloseHandle")
)

// Instance хранит дескриптор мьютекса
type Instance struct {
	handle uintptr
}

// New создаёт новый экземпляр проверки
func New(name string) (*Instance, error) {
	mutexName, err := syscall.UTF16PtrFromString("Local\\" + name)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutex name: %w", err)
	}

	ret, _, errno := procCreateMutexW.Call(
		0,
		0,
		uintptr(unsafe.Pointer(mutexName)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("failed to create mutex: %v", errno)
	}

	if errno == syscall.Errno(183) {
		procCloseHandle.Call(ret)
		return nil, fmt.Errorf("application is already running")
	}

	return &Instance{handle: ret}, nil
}

// Release освобождает мьютекс
func (i *Instance) Release() {
	if i.handle != 0 {
		procCloseHandle.Call(i.handle)
		i.handle = 0
	}
}

// GetMutexName возвращает имя мьютекса
func GetMutexName() string {
	execPath, err := os.Executable()
	if err != nil {
		return "voice-input-go"
	}
	
	base := execPath
	for i := len(base) - 1; i >= 0; i-- {
		if base[i] == '\\' || base[i] == '/' {
			base = base[i+1:]
			break
		}
	}
	
	if len(base) > 4 && base[len(base)-4:] == ".exe" {
		base = base[:len(base)-4]
	}
	
	return base
}
