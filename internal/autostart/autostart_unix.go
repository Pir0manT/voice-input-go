//go:build !windows

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const desktopEntry = `[Desktop Entry]
Type=Application
Name=Voice Input Go
Comment=Voice transcription via Lemonade Server
Exec=%s
Terminal=false
Categories=Utility;
`

func enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	switch runtime.GOOS {
	case "linux":
		return enableLinux(exePath)
	case "darwin":
		return enableMac(exePath)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func disable() error {
	switch runtime.GOOS {
	case "linux":
		return disableLinux()
	case "darwin":
		return disableMac()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func isEnabled() bool {
	switch runtime.GOOS {
	case "linux":
		return isEnabledLinux()
	case "darwin":
		return isEnabledMac()
	default:
		return false
	}
}

// Linux

func enableLinux(exePath string) error {
	dir, err := getAutostartDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create autostart directory: %w", err)
	}

	content := fmt.Sprintf(desktopEntry, exePath)
	path := filepath.Join(dir, "voice-input-go.desktop")
	return os.WriteFile(path, []byte(content), 0755)
}

func disableLinux() error {
	dir, err := getAutostartDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "voice-input-go.desktop")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func isEnabledLinux() bool {
	dir, _ := getAutostartDir()
	if dir == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(dir, "voice-input-go.desktop"))
	return err == nil
}

func getAutostartDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "autostart"), nil
}

// macOS

func enableMac(exePath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.pir0mant.voice-input-go</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>`, exePath)

	path := filepath.Join(dir, "com.pir0mant.voice-input-go.plist")
	return os.WriteFile(path, []byte(plist), 0644)
}

func disableMac() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(home, "Library", "LaunchAgents", "com.pir0mant.voice-input-go.plist")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func isEnabledMac() bool {
	home, _ := os.UserHomeDir()
	if home == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(home, "Library", "LaunchAgents", "com.pir0mant.voice-input-go.plist"))
	return err == nil
}
