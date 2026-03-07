//go:build windows

package console

import "syscall"

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	user32           = syscall.NewLazyDLL("user32.dll")
	getConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	showWindow       = user32.NewProc("ShowWindow")
)

func SetVisible(visible bool) {
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}
	if visible {
		showWindow.Call(hwnd, 5) // SW_SHOW
	} else {
		showWindow.Call(hwnd, 0) // SW_HIDE
	}
}
