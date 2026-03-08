package tray

import (
	"fmt"
	"sync"

	"fyne.io/systray"
	"github.com/Pir0manT/voice-input-go/internal/config"
	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// Status статус приложения
type Status int

const (
	StatusIdle Status = iota
	StatusRecording
	StatusProcessing
)

var (
	mu            sync.Mutex
	currentStatus Status = StatusIdle
	cfg           *config.Config
	appLang       string = "ru"

	statusItem   *systray.MenuItem
	startItem    *systray.MenuItem
	stopItem     *systray.MenuItem
	editorItem   *systray.MenuItem
	settingsItem *systray.MenuItem
	logsItem     *systray.MenuItem
	quitItem     *systray.MenuItem

	onStartCallback    func()
	onStopCallback     func()
	onEditorCallback   func()
	onSettingsCallback func()
	onLogsCallback     func()
	onReadyCallback    func()
)

// Start запускает системный трей
func Start(c *config.Config, callbacks map[string]func(), lang string) error {
	mu.Lock()
	cfg = c
	appLang = lang

	if fn, ok := callbacks["start"]; ok {
		onStartCallback = fn
	}
	if fn, ok := callbacks["stop"]; ok {
		onStopCallback = fn
	}
	if fn, ok := callbacks["editor"]; ok {
		onEditorCallback = fn
	}
	if fn, ok := callbacks["settings"]; ok {
		onSettingsCallback = fn
	}
	if fn, ok := callbacks["logs"]; ok {
		onLogsCallback = fn
	}
	mu.Unlock()

	// Запускаем systray (блокирующий вызов)
	systray.Run(onReady, onExit)

	return nil
}

// Restart обновляет настройки трея без перезапуска
func Restart(c *config.Config, callbacks map[string]func(), lang string) {
	mu.Lock()
	cfg = c
	appLang = lang
	status := currentStatus
	mu.Unlock()

	msg := i18n.Get(lang)

	// Обновляем tooltip
	systray.SetTooltip(msg.TrayTooltipRU)

	// Обновляем все пункты меню
	mu.Lock()
	if startItem != nil {
		startItem.SetTitle("🎤 " + msg.TrayStart)
	}
	if stopItem != nil {
		stopItem.SetTitle("⏹ " + msg.TrayStop)
	}
	if editorItem != nil {
		editorItem.SetTitle("📝 " + msg.TrayEditor)
	}
	if settingsItem != nil {
		settingsItem.SetTitle(msg.TraySettings)
	}
	if logsItem != nil {
		logsItem.SetTitle(msg.TrayLogs)
	}
	if quitItem != nil {
		quitItem.SetTitle("❌ " + msg.TrayQuit)
	}

	// Обновляем statusItem и иконку по текущему статусу
	if statusItem != nil {
		switch status {
		case StatusIdle:
			systray.SetIcon(IconIdle)
			statusItem.SetTitle(msg.StatusIdle)
		case StatusRecording:
			systray.SetIcon(IconRecording)
			statusItem.SetTitle(msg.StatusRecording)
		case StatusProcessing:
			systray.SetIcon(IconProcessing)
			statusItem.SetTitle(msg.StatusProcessing)
		}
	}
	mu.Unlock()
}

// SetOnReady устанавливает callback, который вызывается после инициализации трея.
// На macOS это единственный безопасный момент для регистрации глобальных хоткеев,
// т.к. Carbon API требует работающего NSApplication event loop.
func SetOnReady(fn func()) {
	mu.Lock()
	defer mu.Unlock()
	onReadyCallback = fn
}

// onReady вызывается когда трей готов
func onReady() {
	mu.Lock()
	lang := appLang
	mu.Unlock()

	msg := i18n.Get(lang)

	// Устанавливаем иконку
	if len(IconIdle) > 0 {
		systray.SetIcon(IconIdle)
	}

	systray.SetTitle("") // Пустой title — на Linux текст рядом с иконкой не нужен
	systray.SetTooltip(msg.TrayTooltipRU)

	// Создаём меню
	mu.Lock()
	statusItem = systray.AddMenuItem(msg.StatusIdle, "Current status")
	statusItem.Disable()
	mu.Unlock()

	systray.AddSeparator()

	mu.Lock()
	startItem = systray.AddMenuItem("🎤 "+msg.TrayStart, msg.TrayStart)
	stopItem = systray.AddMenuItem("⏹ "+msg.TrayStop, msg.TrayStop)
	stopItem.Disable()
	mu.Unlock()

	systray.AddSeparator()

	mu.Lock()
	editorItem = systray.AddMenuItem("📝 "+msg.TrayEditor, msg.TrayEditor)
	settingsItem = systray.AddMenuItem(msg.TraySettings, msg.TraySettings)
	logsItem = systray.AddMenuItem(msg.TrayLogs, msg.TrayLogs)
	mu.Unlock()

	systray.AddSeparator()

	mu.Lock()
	quitItem = systray.AddMenuItem("❌ "+msg.TrayQuit, msg.TrayQuit)
	mu.Unlock()

	// Обработчики событий
	go eventLoop()

	// Вызываем onReady callback (регистрация хоткеев и т.д.)
	mu.Lock()
	readyCb := onReadyCallback
	mu.Unlock()
	if readyCb != nil {
		go readyCb()
	}
}

// eventLoop обрабатывает клики по пунктам меню
func eventLoop() {
	for {
		mu.Lock()
		si := startItem
		sti := stopItem
		ei := editorItem
		sei := settingsItem
		li := logsItem
		qi := quitItem
		mu.Unlock()

		select {
		case <-si.ClickedCh:
			mu.Lock()
			cb := onStartCallback
			mu.Unlock()
			if cb != nil {
				cb()
			}
		case <-sti.ClickedCh:
			mu.Lock()
			cb := onStopCallback
			mu.Unlock()
			if cb != nil {
				cb()
			}
		case <-ei.ClickedCh:
			mu.Lock()
			cb := onEditorCallback
			mu.Unlock()
			if cb != nil {
				cb()
			}
		case <-sei.ClickedCh:
			mu.Lock()
			cb := onSettingsCallback
			mu.Unlock()
			if cb != nil {
				cb()
			}
		case <-li.ClickedCh:
			mu.Lock()
			cb := onLogsCallback
			mu.Unlock()
			if cb != nil {
				cb()
			}
		case <-qi.ClickedCh:
			mu.Lock()
			lang := appLang
			mu.Unlock()
			fmt.Println("❌ " + i18n.Get(lang).TrayQuit)
			systray.Quit()
		}
	}
}

// onExit вызывается при выходе
func onExit() {
	mu.Lock()
	lang := appLang
	mu.Unlock()
	msg := i18n.Get(lang)
	fmt.Println("👋 " + msg.Close)
}

// SetStatus устанавливает статус приложения
func SetStatus(status Status) {
	mu.Lock()
	currentStatus = status
	si := statusItem
	sta := startItem
	sto := stopItem
	lang := appLang
	mu.Unlock()

	if si == nil {
		return
	}

	msg := i18n.Get(lang)

	switch status {
	case StatusIdle:
		systray.SetIcon(IconIdle)
		si.SetTitle(msg.StatusIdle)
		sta.Enable()
		sto.Disable()
	case StatusRecording:
		systray.SetIcon(IconRecording)
		si.SetTitle(msg.StatusRecording)
		sta.Disable()
		sto.Enable()
	case StatusProcessing:
		systray.SetIcon(IconProcessing)
		si.SetTitle(msg.StatusProcessing)
		sta.Disable()
		sto.Disable()
	}
}

// GetStatus возвращает текущий статус
func GetStatus() Status {
	mu.Lock()
	defer mu.Unlock()
	return currentStatus
}

// Quit завершает приложение
func Quit() {
	systray.Quit()
}
