package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Pir0manT/voice-input-go/internal/i18n"
)

// Level уровня логирования
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}
	appLogger  *log.Logger
	logFile    *os.File
	currentLevel Level
	mu         sync.Mutex
	currentLang string = "ru"
)

// Init инициализирует логгер
func Init(filename string, level Level, lang string) error {
	mu.Lock()
	defer mu.Unlock()

	currentLang = lang
	msg := i18n.Get(lang)

	// Используем AppData для логов
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		// Фоллбэк на текущую директорию
		appDataDir = "."
	}
	
	logDir := filepath.Join(appDataDir, "voice-input-go", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	
	logPath := filepath.Join(logDir, filepath.Base(filename))

	// Открываем файл
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	currentLevel = level

	// Создаём мультиписатель (консоль + файл)
	mw := io.MultiWriter(os.Stdout, logFile)
	appLogger = log.New(mw, "", log.LstdFlags)

	appLogger.Printf(msg.LoggerInitialized, levelNames[level])
	return nil
}

// Close закрывает лог файл
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		logFile.Close()
	}
}

// writeLog записывает сообщение с указанным уровнем
func writeLog(level Level, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	if level < currentLevel {
		return
	}

	// Добавляем информацию о файле и строке (только для DEBUG)
	if level == DEBUG {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			format = fmt.Sprintf("[%s:%d] %s", filepath.Base(file), line, format)
		}
	}

	msg := fmt.Sprintf(format, v...)
	appLogger.Printf("[%s] %s", levelNames[level], msg)
}

// Debug записывает DEBUG сообщение
func Debug(format string, v ...interface{}) {
	writeLog(DEBUG, format, v...)
}

// Info записывает INFO сообщение
func Info(format string, v ...interface{}) {
	writeLog(INFO, format, v...)
}

// Warn записывает WARN сообщение
func Warn(format string, v ...interface{}) {
	writeLog(WARN, format, v...)
}

// Error записывает ERROR сообщение
func Error(format string, v ...interface{}) {
	writeLog(ERROR, format, v...)
}

// GetLogFilePath возвращает путь к файлу логов
func GetLogFilePath() string {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return ""
	}

	name, err := logFile.Stat()
	if err != nil {
		return ""
	}

	return name.Name()
}

// GetLogs возвращает последние N строк лога
func GetLogs(lines int) ([]string, error) {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return nil, fmt.Errorf("log file not initialized")
	}

	// TODO: Реализовать чтение последних строк
	return []string{}, nil
}

// ClearLogs очищает лог файл
func ClearLogs() error {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return fmt.Errorf("log file not initialized")
	}

	// Пересоздаём файл
	name := logFile.Name()
	logFile.Close()

	err := os.Truncate(name, 0)
	if err != nil {
		return fmt.Errorf("failed to truncate log file: %w", err)
	}

	// Открываем заново
	var err2 error
	logFile, err2 = os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err2 != nil {
		return fmt.Errorf("failed to reopen log file: %w", err2)
	}

	appLogger.Println(i18n.Get(currentLang).LogsCleared)
	return nil
}

// SetLevel устанавливает уровень логирования
func SetLevel(level Level) {
	mu.Lock()
	defer mu.Unlock()
	currentLevel = level
}

// ParseLevel парсит строку уровня логирования
func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// SaveLogs сохраняет логи в указанный файл
func SaveLogs(filename string) error {
	mu.Lock()
	defer mu.Unlock()

	if logFile == nil {
		return fmt.Errorf("log file not initialized")
	}

	// TODO: Реализовать сохранение в другой файл
	return nil
}
