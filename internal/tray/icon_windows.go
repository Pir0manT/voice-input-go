//go:build windows

package tray

import _ "embed"

//go:embed icon_idle.ico
var IconIdle []byte

//go:embed icon_recording.ico
var IconRecording []byte

//go:embed icon_processing.ico
var IconProcessing []byte
