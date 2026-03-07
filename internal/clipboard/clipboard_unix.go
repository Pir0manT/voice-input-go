//go:build !windows

package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func copyPlatform(text string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	case "linux":
		// Пробуем wl-copy (Wayland)
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd := exec.Command("wl-copy")
			cmd.Stdin = strings.NewReader(text)
			if err := cmd.Run(); err == nil {
				return nil
			}
		}
		// Пробуем xclip (X11)
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
		// Пробуем xsel (X11)
		cmd = exec.Command("xsel", "--clipboard", "--input")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func pastePlatform() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		out, err := exec.Command("pbpaste").Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	case "linux":
		// Пробуем wl-paste (Wayland)
		if _, err := exec.LookPath("wl-paste"); err == nil {
			out, err := exec.Command("wl-paste", "--no-newline").Output()
			if err == nil {
				return string(out), nil
			}
		}
		// Пробуем xclip (X11)
		out, err := exec.Command("xclip", "-selection", "clipboard", "-output").Output()
		if err == nil {
			return string(out), nil
		}
		// Пробуем xsel (X11)
		out, err = exec.Command("xsel", "--clipboard", "--output").Output()
		if err == nil {
			return string(out), nil
		}
		return "", fmt.Errorf("failed to paste: wl-paste, xclip and xsel not available")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
