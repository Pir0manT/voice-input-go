package notify

import (
	"fmt"
	"os/exec"
	"runtime"
)

// PlaySound воспроизводит системный звук уведомления
func PlaySound() error {
	switch runtime.GOOS {
	case "windows":
		return playSoundWindows()
	case "darwin":
		return playSoundMac()
	case "linux":
		return playSoundLinux()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func playSoundWindows() error {
	cmd := exec.Command("powershell", "-NoProfile", "-command",
		"[System.Media.SystemSounds]::Beep.Play()")
	return cmd.Run()
}

func playSoundMac() error {
	cmd := exec.Command("afplay", "/System/Library/Sounds/Glass.aiff")
	return cmd.Run()
}

func playSoundLinux() error {
	cmd := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga")
	if err := cmd.Run(); err == nil {
		return nil
	}
	cmd = exec.Command("aplay", "/usr/share/sounds/alsa/Front_Center.wav")
	return cmd.Run()
}
