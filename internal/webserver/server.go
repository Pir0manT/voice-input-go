package webserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Pir0manT/voice-input-go/internal/config"
	"github.com/Pir0manT/voice-input-go/internal/i18n"
	"github.com/skratchdot/open-golang/open"
)

// Server веб сервер
type Server struct {
	port         int
	config       *config.Config
	mu           sync.RWMutex
	server       *http.Server
	messages     *i18n.Messages
	onLangChange func(string) // Callback при смене языка
}

// ConfigRequest запрос на обновление конфига
type ConfigRequest struct {
	Hotkeys       config.HotkeysConfig       `json:"hotkeys"`
	Lemonade      config.LemonadeConfig      `json:"lemonade"`
	Notifications config.NotificationsConfig `json:"notifications"`
	Autostart     bool                       `json:"autostart"`
	Logging       config.LoggingConfig       `json:"logging"`
	AppLanguage   string                     `json:"appLanguage"`
	Restart       bool                       `json:"restart"` // Флаг перезапуска трея
}

// templateData данные для шаблона
type templateData struct {
	Lang       string
	Messages   *i18n.Messages
	Config     *config.Config
}

// NewServer создаёт новый сервер
func NewServer(cfg *config.Config) *Server {
	return &Server{
		port:     8080,
		config:   cfg,
		messages: i18n.Get(cfg.AppLanguage),
	}
}

// SetLangChangeCallback устанавливает callback для смены языка
func (s *Server) SetLangChangeCallback(fn func(string)) {
	s.onLangChange = fn
}

// Start запускает сервер
func (s *Server) Start() error {
	// Находим свободный порт
	for !isPortFree(s.port) {
		s.port++
	}

	mux := http.NewServeMux()

	// Статические файлы (HTML, CSS, JS)
	mux.HandleFunc("/", s.handleStatic)

	// API endpoints
	mux.HandleFunc("/api/config", s.handleConfig)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", s.port),
		Handler: s.corsMiddleware(mux),
	}

	fmt.Printf(s.messages.WebServerStarted+"\n", s.port)
	fmt.Printf(s.messages.WebServerSettings+"\n", s.port)

	return s.server.ListenAndServe()
}

// Stop останавливает сервер
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

// OpenSettings открывает настройки в браузере
func (s *Server) OpenSettings() error {
	url := fmt.Sprintf("http://localhost:%d/settings", s.port)
	return open.Run(url)
}

// handleStatic обрабатывает статические файлы
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Если корень или /settings, рендерим settings.html с template
	if path == "/" || path == "/settings" {
		s.renderSettingsPage(w)
		return
	}

	// Для других файлов пытаемся загрузить из embed
	// path начинается с /, убираем его
	filePath := path
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	data, err := GetFile(filePath)
	if err != nil {
		http.Error(w, s.messages.WebServerFileNotFound, http.StatusNotFound)
		return
	}

	// Определяем Content-Type по расширению
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	}

	w.Write(data)
}

// renderSettingsPage рендерит страницу настроек
func (s *Server) renderSettingsPage(w http.ResponseWriter) {
	data := templateData{
		Lang:     s.config.AppLanguage,
		Messages: s.messages,
		Config:   s.config,
	}

	html, err := GetFile("settings.html")
	if err != nil {
		http.Error(w, s.messages.WebServerFileNotFound, http.StatusNotFound)
		return
	}

	// Заменяем плейсхолдеры в HTML
	htmlStr := string(html)
	htmlStr = s.applyTranslations(htmlStr, data)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(htmlStr))
}

// applyTranslations применяет переводы к HTML
func (s *Server) applyTranslations(html string, data templateData) string {
	// Заменяем плейсхолдеры вида {{.Messages.Save}}
	replacements := map[string]string{
		"{{.Messages.Save}}":              data.Messages.Save,
		"{{.Messages.Cancel}}":            data.Messages.Cancel,
		"{{.Messages.SettingsTitle}}":     data.Messages.SettingsTitle,
		"{{.Messages.TabGeneral}}":        data.Messages.TabGeneral,
		"{{.Messages.TabNotifications}}":  data.Messages.TabNotifications,
		"{{.Messages.TabLogs}}":           data.Messages.TabLogs,
		"{{.Messages.SectionHotkeys}}":    data.Messages.SectionHotkeys,
		"{{.Messages.LabelHotkeyStart}}":  data.Messages.LabelHotkeyStart,
		"{{.Messages.LabelHotkeyStop}}":   data.Messages.LabelHotkeyStop,
		"{{.Messages.LabelHotkeyEditor}}": data.Messages.LabelHotkeyEditor,
		"{{.Messages.SectionLemonade}}":   data.Messages.SectionLemonade,
		"{{.Messages.LabelURL}}":          data.Messages.LabelURL,
		"{{.Messages.LabelModel}}":        data.Messages.LabelModel,
		"{{.Messages.LabelLanguage}}":     data.Messages.LabelLanguage,
		"{{.Messages.SectionAutostart}}":  data.Messages.SectionAutostart,
		"{{.Messages.CheckboxAutostart}}": data.Messages.CheckboxAutostart,
		"{{.Messages.HintAutostart}}":     data.Messages.HintAutostart,
		"{{.Messages.SectionAppLanguage}}": data.Messages.SectionAppLanguage,
		"{{.Messages.LabelAppLanguage}}":   data.Messages.LabelAppLanguage,
		"{{.Messages.SectionNotifications}}": data.Messages.SectionNotifications,
		"{{.Messages.CheckboxSound}}":      data.Messages.CheckboxSound,
		"{{.Messages.CheckboxToast}}":      data.Messages.CheckboxToast,
		"{{.Messages.SectionLogs}}":        data.Messages.SectionLogs,
		"{{.Messages.CheckboxLogging}}":    data.Messages.CheckboxLogging,
		"{{.Messages.LabelLoggingLevel}}":  data.Messages.LabelLoggingLevel,
		"{{.Messages.LabelLastLogs}}":      data.Messages.LabelLastLogs,
		"{{.Messages.BtnViewLogs}}":        data.Messages.BtnViewLogs,
		"{{.Messages.BtnClearLogs}}":       data.Messages.BtnClearLogs,
		"{{.Messages.BtnSaveLogs}}":        data.Messages.BtnSaveLogs,
		"{{.Messages.LangRussian}}":        data.Messages.LangRussian,
		"{{.Messages.LangEnglish}}":        data.Messages.LangEnglish,
		"{{.Messages.ModelWhisperTurbo}}":  data.Messages.ModelWhisperTurbo,
	}

	for placeholder, value := range replacements {
		html = strings.ReplaceAll(html, placeholder, value)
	}

	// Устанавливаем lang атрибут
	html = strings.ReplaceAll(html, `<html lang="ru">`, fmt.Sprintf(`<html lang="%s">`, data.Lang))
	html = strings.ReplaceAll(html, `<html lang="en">`, fmt.Sprintf(`<html lang="%s">`, data.Lang))

	return html
}

// handleConfig обрабатывает GET и POST запросы к конфиг
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.handleGetConfig(w, r)
	} else if r.Method == http.MethodPost {
		s.handleSaveConfig(w, r)
	} else {
		http.Error(w, s.messages.WebServerMethodNotAllowed, http.StatusMethodNotAllowed)
	}
}

// handleGetConfig возвращает текущий конфиг
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.config)
}

// handleSaveConfig сохраняет конфиг
func (s *Server) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	var req ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем смену языка
	oldLang := s.config.AppLanguage
	langChanged := req.AppLanguage != oldLang

	// Обновляем конфиг
	s.config.Hotkeys = req.Hotkeys
	s.config.Lemonade = req.Lemonade
	s.config.Notifications = req.Notifications
	s.config.Autostart = req.Autostart
	s.config.Logging = req.Logging
	s.config.AppLanguage = req.AppLanguage

	// Сохраняем на диск
	if err := config.Save(s.config); err != nil {
		http.Error(w, s.messages.WebServerErrorSaving, http.StatusInternalServerError)
		return
	}

	// Если язык изменился, вызываем callback
	if langChanged && s.onLangChange != nil {
		go s.onLangChange(req.AppLanguage)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// corsMiddleware добавляет CORS заголовки для localhost
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+fmt.Sprint(s.port))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isPortFree проверяет свободен ли порт
func isPortFree(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}
