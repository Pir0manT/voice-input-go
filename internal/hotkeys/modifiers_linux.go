//go:build linux

package hotkeys

import "golang.design/x/hotkey"

var (
	modAlt   = hotkey.Mod1 // X11: Mod1 = Alt
	modSuper = hotkey.Mod4 // X11: Mod4 = Super/Win
)
