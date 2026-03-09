package lemonade

import (
	"github.com/Pir0manT/voice-input-go/internal/transcriber"
)

// TranscriberAdapter адаптер, оборачивающий Client в интерфейс transcriber.Transcriber
type TranscriberAdapter struct {
	client      *Client
	model       string
	language    string
	prompt      string
	temperature float64
}

// NewTranscriberAdapter создаёт адаптер для Lemonade клиента
func NewTranscriberAdapter(client *Client, model, language, prompt string, temperature float64) *TranscriberAdapter {
	return &TranscriberAdapter{
		client:      client,
		model:       model,
		language:    language,
		prompt:      prompt,
		temperature: temperature,
	}
}

// SetModel обновляет модель
func (a *TranscriberAdapter) SetModel(model string) {
	a.model = model
}

// SetLanguage обновляет язык
func (a *TranscriberAdapter) SetLanguage(language string) {
	a.language = language
}

// SetPrompt обновляет prompt
func (a *TranscriberAdapter) SetPrompt(prompt string) {
	a.prompt = prompt
}

// SetTemperature обновляет температуру
func (a *TranscriberAdapter) SetTemperature(temperature float64) {
	a.temperature = temperature
}

// TranscribeWithStats реализует transcriber.Transcriber
func (a *TranscriberAdapter) TranscribeWithStats(audioPath string) (*transcriber.Result, error) {
	result, err := a.client.TranscribeWithStats(audioPath, a.model, a.language, a.prompt, a.temperature)
	if err != nil {
		return nil, err
	}

	return &transcriber.Result{
		Text:        result.Text,
		Duration:    result.Duration,
		ProcessTime: result.ProcessTime,
		Speed:       result.Speed,
		Backend:     result.Backend,
	}, nil
}

// CheckHealth реализует transcriber.Transcriber
func (a *TranscriberAdapter) CheckHealth() (bool, error) {
	return a.client.CheckHealth()
}

// GetClient возвращает нижележащий Lemonade клиент (для управления моделями и т.д.)
func (a *TranscriberAdapter) GetClient() *Client {
	return a.client
}
