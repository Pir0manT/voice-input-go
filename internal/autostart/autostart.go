package autostart

// Enable включает автозапуск приложения при старте системы
func Enable() error {
	return enable()
}

// Disable отключает автозапуск
func Disable() error {
	return disable()
}

// IsEnabled проверяет включен ли автозапуск
func IsEnabled() bool {
	return isEnabled()
}
