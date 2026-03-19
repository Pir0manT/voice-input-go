package settings

import (
	"fmt"
	"sync"

	"github.com/Pir0manT/voice-input-go/internal/config"
	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// Settings менеджер окна настроек
type Settings struct {
	mu             sync.Mutex
	running        bool
	lang           string
	config         *config.Config
	logPath        string
	onConfigChange func(*config.Config)
	onOpen         func() // вызывается при открытии окна (отключить хоткеи)
	onClose        func() // вызывается при закрытии окна (включить хоткеи)
}

// New создаёт менеджер настроек
func New(cfg *config.Config, lang, logPath string) *Settings {
	return &Settings{
		config:  cfg,
		lang:    lang,
		logPath: logPath,
	}
}

// SetOnConfigChange устанавливает callback при сохранении настроек
func (s *Settings) SetOnConfigChange(fn func(*config.Config)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onConfigChange = fn
}

// SetOnOpen устанавливает callback при открытии окна настроек
func (s *Settings) SetOnOpen(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onOpen = fn
}

// SetOnClose устанавливает callback при закрытии окна настроек
func (s *Settings) SetOnClose(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onClose = fn
}

// SetLanguage обновляет язык
func (s *Settings) SetLanguage(lang string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lang = lang
}

// SetConfig обновляет конфиг
func (s *Settings) SetConfig(cfg *config.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = cfg
}

// IsRunning проверяет открыто ли окно
func (s *Settings) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Show открывает окно настроек
func (s *Settings) Show() {
	s.show(0)
}

// ShowLogs открывает окно настроек на вкладке логов
func (s *Settings) ShowLogs() {
	s.show(-1) // -1 = последняя вкладка
}

func (s *Settings) show(initialTab int) {
	s.mu.Lock()
	msg := i18n.Get(s.lang)

	if s.running {
		s.mu.Unlock()
		fmt.Println(msg.SettingsAlreadyOpen)
		return
	}

	s.running = true
	input := SettingsInput{
		Config:     s.config,
		Lang:       s.lang,
		LogPath:    s.logPath,
		InitialTab: initialTab,
	}
	onOpen := s.onOpen
	onClose := s.onClose
	s.mu.Unlock()

	fmt.Println(msg.SettingsOpening)

	// Отключаем глобальные хоткеи чтобы они не перехватывали нажатия в окне настроек
	if onOpen != nil {
		onOpen()
	}

	go func() {
		defer func() {
			// Включаем хоткеи обратно при закрытии окна
			if onClose != nil {
				onClose()
			}
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}()

		output, err := launchSettingsWindow(input)
		if err != nil {
			fmt.Printf(msg.SettingsProcessError+"\n", err)
			return
		}

		if output == nil {
			return
		}

		// Сохраняем конфиг на диск
		if err := config.Save(output.Config); err != nil {
			fmt.Printf(msg.MsgSettingsError+"%v\n", err)
			return
		}

		fmt.Println(msg.MsgSettingsSaved)

		// Вызываем callback
		s.mu.Lock()
		cb := s.onConfigChange
		s.mu.Unlock()

		if cb != nil {
			cb(output.Config)
		}
	}()
}
