//go:build windows

package notify

import (
	"github.com/go-toast/toast"
)

// ShowToast показывает всплывающее уведомление Windows через WinRT Toast API.
func ShowToast(title, message string) error {
	notification := toast.Notification{
		AppID:   "Voice Input Go",
		Title:   title,
		Message: message,
	}
	return notification.Push()
}
