package fastflowlm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Pir0manT/voice-input-go/internal/transcriber"
)

// Client клиент FastFlowLM API (OpenAI-совместимый)
type Client struct {
	url        string
	model      string
	language   string
	prompt     string
	httpClient *http.Client
}

// transcribeResponse ответ /v1/audio/transcriptions
type transcribeResponse struct {
	Text string `json:"text"`
}

// NewClient создаёт новый клиент FastFlowLM
func NewClient(url, model, language, prompt string) *Client {
	return &Client{
		url:      url,
		model:    model,
		language: language,
		prompt:   prompt,
		httpClient: &http.Client{
			Timeout: 2 * time.Minute,
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

// SetModel обновляет модель
func (c *Client) SetModel(model string) {
	c.model = model
}

// TranscribeWithStats отправляет аудио на транскрибацию через POST /v1/audio/transcriptions
func (c *Client) TranscribeWithStats(audioPath string) (*transcriber.Result, error) {
	startTime := time.Now()

	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("audio file not found: %s", audioPath)
	}

	// Получаем длительность аудио из WAV файла
	duration, err := getWavDuration(audioPath)
	if err != nil {
		duration = 0
	}

	// Создаём multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем аудио файл — поле "file" (OpenAI-совместимый API)
	file, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Добавляем модель
	if err := writer.WriteField("model", c.model); err != nil {
		return nil, fmt.Errorf("failed to write model field: %w", err)
	}

	// Добавляем язык (опционально)
	if c.language != "" {
		if err := writer.WriteField("language", c.language); err != nil {
			return nil, fmt.Errorf("failed to write language field: %w", err)
		}
	}

	// Добавляем промпт (опционально)
	if c.prompt != "" {
		if err := writer.WriteField("prompt", c.prompt); err != nil {
			return nil, fmt.Errorf("failed to write prompt field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Создаём HTTP запрос
	reqURL := c.url + "/v1/audio/transcriptions"
	req, err := http.NewRequest("POST", reqURL, body)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error (%d): %s", resp.StatusCode, string(respBody))
	}

	// Парсим JSON ответ
	var result transcribeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	processTime := time.Since(startTime).Seconds()
	speed := 0.0
	if processTime > 0 && duration > 0 {
		speed = duration / processTime
	}

	return &transcriber.Result{
		Text:        result.Text,
		Duration:    duration,
		ProcessTime: processTime,
		Speed:       speed,
		Backend:     "fastflowlm",
	}, nil
}

// CheckHealth проверяет доступность сервера (GET /v1/models)
func (c *Client) CheckHealth() (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(c.url + "/v1/models")
	if err != nil {
		return false, fmt.Errorf("server unavailable: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

// getWavDuration возвращает длительность WAV файла в секундах
func getWavDuration(filePath string) (float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	header := make([]byte, 44)
	if _, err := io.ReadFull(file, header); err != nil {
		return 0, err
	}

	sampleRate := uint32(header[24]) | uint32(header[25])<<8 | uint32(header[26])<<16 | uint32(header[27])<<24
	dataSize := uint32(header[40]) | uint32(header[41])<<8 | uint32(header[42])<<16 | uint32(header[43])<<24
	bitsPerSample := uint32(header[34]) | uint32(header[35])<<8
	numChannels := uint32(header[22]) | uint32(header[23])<<8

	bytesPerSample := bitsPerSample / 8
	numSamples := dataSize / bytesPerSample
	duration := float64(numSamples) / float64(sampleRate*numChannels)

	return duration, nil
}
