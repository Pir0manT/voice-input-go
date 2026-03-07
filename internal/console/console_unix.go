//go:build !windows

package console

import (
	"os"
	"os/exec"
	"sync"

	"github.com/Pir0manT/voice-input-go/internal/config"
)

var (
	termProcess *os.Process
	mu          sync.Mutex
)

// SetVisible на Linux: открывает/закрывает окно терминала с логами.
// При автозапуске процесс не привязан к терминалу, поэтому
// для просмотра вывода открываем отдельное окно с tail -f на лог-файл.
func SetVisible(visible bool) {
	mu.Lock()
	defer mu.Unlock()

	if visible {
		// Если терминал уже открыт — не открываем второй
		if termProcess != nil {
			return
		}

		logPath, err := config.GetLogFilePath()
		if err != nil {
			return
		}

		// Пробуем разные эмуляторы терминалов
		cmd := findTerminal(logPath)
		if cmd == nil {
			return
		}

		if err := cmd.Start(); err != nil {
			return
		}
		termProcess = cmd.Process
	} else {
		if termProcess != nil {
			termProcess.Kill()
			termProcess.Wait()
			termProcess = nil
		}
	}
}

// findTerminal находит доступный эмулятор терминала и создаёт команду
func findTerminal(logPath string) *exec.Cmd {
	title := "Voice Input Go — Logs"

	// gnome-terminal (GNOME)
	if path, err := exec.LookPath("gnome-terminal"); err == nil {
		return exec.Command(path, "--title", title, "--", "tail", "-f", logPath)
	}

	// konsole (KDE)
	if path, err := exec.LookPath("konsole"); err == nil {
		return exec.Command(path, "--title", title, "-e", "tail", "-f", logPath)
	}

	// xfce4-terminal (XFCE)
	if path, err := exec.LookPath("xfce4-terminal"); err == nil {
		return exec.Command(path, "--title", title, "-e", "tail -f "+logPath)
	}

	// xterm (fallback)
	if path, err := exec.LookPath("xterm"); err == nil {
		return exec.Command(path, "-title", title, "-e", "tail", "-f", logPath)
	}

	return nil
}
