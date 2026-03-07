//go:build darwin || linux
// +build darwin linux

package singleton

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Instance хранит файл блокировки
type Instance struct {
	file *os.File
	path string
}

// New создаёт новый экземпляр проверки
func New(name string) (*Instance, error) {
	// Используем файл блокировки в /tmp или /var/tmp
	tmpDir := os.TempDir()
	lockPath := filepath.Join(tmpDir, fmt.Sprintf("%s.lock", name))

	// Открываем или создаём файл блокировки
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create lock file: %w", err)
	}

	// Пробуем заблокировать файл (неблокирующая попытка)
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("application is already running")
	}

	return &Instance{file: file, path: lockPath}, nil
}

// Release освобождает блокировку
func (i *Instance) Release() {
	if i.file != nil {
		syscall.Flock(int(i.file.Fd()), syscall.LOCK_UN)
		i.file.Close()
		os.Remove(i.path)
	}
}

// GetMutexName возвращает имя для блокировки
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
	
	// Убираем расширение .exe если есть
	if len(base) > 4 && base[len(base)-4:] == ".exe" {
		base = base[:len(base)-4]
	}
	
	return base
}
