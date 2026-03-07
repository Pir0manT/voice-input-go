//go:build !windows

package notify

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ShowToast показывает всплывающее уведомление (macOS/Linux)
func ShowToast(title, message string) error {
	switch runtime.GOOS {
	case "darwin":
		return showToastMacOS(title, message)
	case "linux":
		return showToastLinux(title, message)
	default:
		return fmt.Errorf("toast notifications not supported on %s", runtime.GOOS)
	}
}

func showToastMacOS(title, message string) error {
	// osascript передаёт строки через -e, кириллица работает нативно.
	// Экранируем обратные слэши и кавычки для AppleScript строк.
	title = strings.ReplaceAll(title, `\`, `\\`)
	title = strings.ReplaceAll(title, `"`, `\"`)
	message = strings.ReplaceAll(message, `\`, `\\`)
	message = strings.ReplaceAll(message, `"`, `\"`)

	script := fmt.Sprintf(
		`display notification "%s" with title "%s"`,
		message, title,
	)
	return exec.Command("osascript", "-e", script).Run()
}

func showToastLinux(title, message string) error {
	// notify-send (libnotify) — стандарт для GNOME/KDE/XFCE
	if path, err := exec.LookPath("notify-send"); err == nil {
		return exec.Command(path, "--app-name=Voice Input Go", title, message).Run()
	}

	// kdialog — KDE
	if path, err := exec.LookPath("kdialog"); err == nil {
		return exec.Command(path, "--passivepopup", message, "5", "--title", title).Run()
	}

	// zenity — GNOME fallback
	if path, err := exec.LookPath("zenity"); err == nil {
		return exec.Command(path, "--notification", "--text="+title+": "+message).Run()
	}

	return fmt.Errorf("no notification tool found (tried: notify-send, kdialog, zenity)")
}
