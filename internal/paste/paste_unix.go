//go:build !windows

package paste

import (
	"os/exec"
	"runtime"
	"time"
)

// SimulateCtrlV эмулирует нажатие Ctrl+V
func SimulateCtrlV() {
	// Небольшая задержка чтобы буфер обмена успел обновиться
	time.Sleep(100 * time.Millisecond)

	switch runtime.GOOS {
	case "linux":
		// xdotool — стандартный инструмент для эмуляции ввода в X11/XWayland
		if path, err := exec.LookPath("xdotool"); err == nil {
			exec.Command(path, "key", "ctrl+v").Run()
		}
	case "darwin":
		// osascript для macOS
		exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`).Run()
	}
}
