package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Backend тип бэкенда транскрибации
const (
	BackendLemonade  = "lemonade"
	BackendWhisperAPI = "whisper-api"
)

// Config структура конфигурации приложения
type Config struct {
	Hotkeys       HotkeysConfig       `json:"hotkeys"`
	Backend       string              `json:"backend"`       // "lemonade" или "whisper-api"
	Lemonade      LemonadeConfig      `json:"lemonade"`
	WhisperAPI    WhisperAPIConfig    `json:"whisperApi"`
	Notifications NotificationsConfig `json:"notifications"`
	Autostart     bool                `json:"autostart"`
	AutoPaste     bool                `json:"autoPaste"`
	Logging       LoggingConfig       `json:"logging"`
	AppLanguage   string              `json:"appLanguage"`   // Язык интерфейса: "ru" или "en"
	HistorySize   int                 `json:"historySize"`   // Количество записей в истории (по умолчанию 20)
	ShowConsole   bool                `json:"showConsole"`   // Показывать консольное окно (default: false)
}

// HotkeysConfig конфигурация горячих клавиш
type HotkeysConfig struct {
	Start  string `json:"start"`
	Stop   string `json:"stop"`
	Editor string `json:"editor"`
}

// LemonadeConfig конфигурация Lemonade Server
type LemonadeConfig struct {
	URL         string  `json:"url"`
	Model       string  `json:"model"`
	Language    string  `json:"language"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
}

// WhisperAPIConfig конфигурация внешнего Whisper API (whisper-asr-webservice и совместимые)
type WhisperAPIConfig struct {
	URL      string `json:"url"`      // Базовый URL (например http://192.168.1.50:9000)
	Language string `json:"language"` // Код языка (ru, en, de...) — пустой = автоопределение
	Prompt   string `json:"prompt"`   // initial_prompt — подсказка для контекста
}

// NotificationsConfig конфигурация уведомлений
type NotificationsConfig struct {
	Sound         bool `json:"sound"`
	Toast         bool `json:"toast"`
	SoundOnRecord bool `json:"soundOnRecord"` // Звук при начале записи
}

// LoggingConfig конфигурация логирования
type LoggingConfig struct {
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"`
}

// Default возвращает конфигурацию по умолчанию
func Default() *Config {
	return &Config{
		Hotkeys: HotkeysConfig{
			Start:  "alt+r",
			Stop:   "alt+s",
			Editor: "alt+e",
		},
		Backend: BackendLemonade,
		Lemonade: LemonadeConfig{
			URL:         "http://localhost:8000",
			Model:       "Whisper-Large-v3-Turbo",
			Language:    "ru",
			Prompt:      "Привет! Сегодня работаем с Claude Code, GitHub и voice-input. Используем Go, Fyne, PortAudio. Настройки хранятся в AppData. Точки, запятые — всё на месте.",
			Temperature: 0.2,
		},
		WhisperAPI: WhisperAPIConfig{
			URL:      "http://localhost:9000",
			Language: "ru",
			Prompt:   "",
		},
		Notifications: NotificationsConfig{
			Sound: true,
			Toast: true,
		},
		Autostart: false,
		AutoPaste: false,
		Logging: LoggingConfig{
			Enabled: true,
			Level:   "info",
		},
		AppLanguage: "ru", // Русский по умолчанию
		HistorySize: 20,
		ShowConsole: false,
	}
}

// Load загружает конфигурацию из файла
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Проверяем существование файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Файл не существует, создаём дефолтный
		cfg := Default()
		if err := Save(cfg); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		return cfg, nil
	}

	// Читаем файл
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Нормализация: если Backend не задан, ставим lemonade
	if cfg.Backend == "" {
		cfg.Backend = BackendLemonade
	}

	// Нормализация: если HistorySize не задан, ставим дефолт
	if cfg.HistorySize <= 0 {
		cfg.HistorySize = 20
	}

	return &cfg, nil
}

// Save сохраняет конфигурацию в файл
func Save(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Создаём директорию
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Сериализуем
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Записываем файл
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath возвращает путь к файлу конфигурации
func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// getConfigDir возвращает директорию конфигурации
func getConfigDir() (string, error) {
	// Получаем AppData директорию
	appData, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir: %w", err)
	}

	return filepath.Join(appData, "voice-input-go"), nil
}

// GetLogFilePath возвращает путь к файлу логов
func GetLogFilePath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "logs", "voice-input-go.log"), nil
}
