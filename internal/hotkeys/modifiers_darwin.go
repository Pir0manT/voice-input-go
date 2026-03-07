//go:build darwin

package hotkeys

import "golang.design/x/hotkey"

var (
	modAlt   = hotkey.ModOption // macOS: Option = Alt
	modSuper = hotkey.ModCmd    // macOS: Cmd = Super
)
