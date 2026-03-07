//go:build windows

package clipboard

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	ntdll            = syscall.NewLazyDLL("ntdll.dll")
	openClipboard    = user32.NewProc("OpenClipboard")
	closeClipboard   = user32.NewProc("CloseClipboard")
	emptyClipboard   = user32.NewProc("EmptyClipboard")
	setClipboardData = user32.NewProc("SetClipboardData")
	getClipboardData = user32.NewProc("GetClipboardData")
	globalAlloc      = kernel32.NewProc("GlobalAlloc")
	globalFree       = kernel32.NewProc("GlobalFree")
	globalLock       = kernel32.NewProc("GlobalLock")
	globalUnlock     = kernel32.NewProc("GlobalUnlock")
	globalSize       = kernel32.NewProc("GlobalSize")
	rtlMoveMemory    = ntdll.NewProc("RtlMoveMemory")
)

const (
	cfUnicodeText = 13
	gmemMoveable  = 0x0002
)

func copyPlatform(text string) error {
	// Конвертируем UTF-8 в UTF-16
	utf16Text := utf16.Encode([]rune(text))
	utf16Text = append(utf16Text, 0) // null-terminator

	size := len(utf16Text) * 2

	// Выделяем глобальную память
	hMem, _, err := globalAlloc.Call(gmemMoveable, uintptr(size))
	if hMem == 0 {
		return fmt.Errorf("GlobalAlloc failed: %v", err)
	}

	// Блокируем и копируем
	ptr, _, err := globalLock.Call(hMem)
	if ptr == 0 {
		globalFree.Call(hMem)
		return fmt.Errorf("GlobalLock failed: %v", err)
	}

	// Копируем через RtlMoveMemory — безопасно для go vet
	rtlMoveMemory.Call(ptr, uintptr(unsafe.Pointer(&utf16Text[0])), uintptr(size))

	globalUnlock.Call(hMem)

	// Открываем буфер обмена
	ret, _, err := openClipboard.Call(0)
	if ret == 0 {
		globalFree.Call(hMem)
		return fmt.Errorf("OpenClipboard failed: %v", err)
	}
	defer closeClipboard.Call()

	emptyClipboard.Call()

	ret, _, err = setClipboardData.Call(cfUnicodeText, hMem)
	if ret == 0 {
		return fmt.Errorf("SetClipboardData failed: %v", err)
	}
	// После SetClipboardData система владеет памятью — не вызываем GlobalFree

	return nil
}

func pastePlatform() (string, error) {
	ret, _, err := openClipboard.Call(0)
	if ret == 0 {
		return "", fmt.Errorf("OpenClipboard failed: %v", err)
	}
	defer closeClipboard.Call()

	hMem, _, _ := getClipboardData.Call(cfUnicodeText)
	if hMem == 0 {
		return "", nil
	}

	memSize, _, _ := globalSize.Call(hMem)
	if memSize == 0 {
		return "", nil
	}

	ptr, _, err := globalLock.Call(hMem)
	if ptr == 0 {
		return "", fmt.Errorf("GlobalLock failed: %v", err)
	}
	defer globalUnlock.Call(hMem)

	// Копируем данные из глобальной памяти в Go slice через RtlMoveMemory
	maxChars := int(memSize) / 2
	buf := make([]uint16, maxChars)
	rtlMoveMemory.Call(uintptr(unsafe.Pointer(&buf[0])), ptr, memSize)

	// Находим null-terminator
	var utf16Chars []uint16
	for _, ch := range buf {
		if ch == 0 {
			break
		}
		utf16Chars = append(utf16Chars, ch)
	}

	return string(utf16.Decode(utf16Chars)), nil
}
