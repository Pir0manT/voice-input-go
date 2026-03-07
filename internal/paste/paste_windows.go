//go:build windows

package paste

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	sendInput  = user32.NewProc("SendInput")
)

const (
	inputKeyboard = 1
	keyEventUp    = 0x0002
	vkControl     = 0x11
	vkV           = 0x56
)

type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	inputType uint32
	ki        keyboardInput
	padding   [8]byte
}

// SimulateCtrlV эмулирует нажатие Ctrl+V через SendInput
func SimulateCtrlV() {
	// Небольшая задержка чтобы буфер обмена успел обновиться
	time.Sleep(100 * time.Millisecond)

	inputs := []input{
		{inputType: inputKeyboard, ki: keyboardInput{wVk: vkControl}},          // Ctrl down
		{inputType: inputKeyboard, ki: keyboardInput{wVk: vkV}},               // V down
		{inputType: inputKeyboard, ki: keyboardInput{wVk: vkV, dwFlags: keyEventUp}},       // V up
		{inputType: inputKeyboard, ki: keyboardInput{wVk: vkControl, dwFlags: keyEventUp}},  // Ctrl up
	}

	sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)
}
