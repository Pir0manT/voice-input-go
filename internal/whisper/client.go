package whisper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/Pir0manT/voice-input-go/internal/transcriber"
)

// Client клиент Whisper ASR Webservice (whisper-asr-webservice и совместимые)
type Client struct {
	url        string
	language   string
	prompt     string
	httpClient *http.Client
}

// asrResponse ответ /asr?output=json
type asrResponse struct {
	Text     string  `json:"text"`
	Language string  `json:"language"`
	Duration float64 `json:"duration"`
}

// NewClient создаёт новый клиент Whisper API
func NewClient(url, language, prompt string) *Client {
	return &Client{
		url:      url,
		language: language,
		prompt:   prompt,
		httpClient: &http.Client{
			Timeout: 2 * time.Minute, // Whisper large-v3 может обрабатывать долго
		},
	}
}

// SetLanguage обновляет язык
func (c *Client) SetLanguage(language string) {
	c.language = language
}

// SetPrompt обновляет prompt
func (c *Client) SetPrompt(prompt string) {
	c.prompt = prompt
}

// TranscribeWithStats отправляет аудио на транскрибацию через POST /asr
func (c *Client) TranscribeWithStats(audioPath string) (*transcriber.Result, error) {
	startTime := time.Now()

	// Проверяем существование файла
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("audio file not found: %s", audioPath)
	}

	// Создаём multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем аудио файл — поле ОБЯЗАТЕЛЬНО "audio_file" (не "file"!)
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("audio_file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Собираем URL с query-параметрами (через net/url для корректного кодирования)
	reqURL, err := url.Parse(c.url + "/asr")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	q := reqURL.Query()
	q.Set("output", "json")
	q.Set("encode", "true")
	q.Set("task", "transcribe")
	if c.language != "" {
		q.Set("language", c.language)
	}
	if c.prompt != "" {
		q.Set("initial_prompt", c.prompt)
	}
	reqURL.RawQuery = q.Encode()

	// Создаём HTTP запрос
	req, err := http.NewRequest("POST", reqURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error (%d): %s", resp.StatusCode, string(respBody))
	}

	// Парсим JSON ответ
	var result asrResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	processTime := time.Since(startTime).Seconds()
	duration := result.Duration
	speed := 0.0
	if processTime > 0 && duration > 0 {
		speed = duration / processTime
	}

	return &transcriber.Result{
		Text:        result.Text,
		Duration:    duration,
		ProcessTime: processTime,
		Speed:       speed,
		Backend:     "whisper-api",
	}, nil
}

// CheckHealth проверяет доступность сервера (GET /)
func (c *Client) CheckHealth() (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(c.url + "/")
	if err != nil {
		return false, fmt.Errorf("server unavailable: %w", err)
	}
	defer resp.Body.Close()

	// Whisper ASR webservice отвечает 200 на GET /
	return resp.StatusCode == http.StatusOK, nil
}
