//go:build windows

package hotkeys

import "golang.design/x/hotkey"

var (
	modAlt   = hotkey.ModAlt
	modSuper = hotkey.ModWin
)
