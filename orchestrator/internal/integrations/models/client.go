package models

import (
	"fmt"
	"net/http"
	"orchestrator/internal/domain"
	"time"
)

// Определение базового URL для каждого сервиса.
const (
	DeepfakeDetectURL        = "http://localhost:8000/deepfake-detect/"
	OpenClosedEyeDetectURL   = "http://localhost:8000/open-closed-eye-detect/"
	AudioFakeDetectionURL    = "http://localhost:8000/audio-fake-detection/"
	LipsMovementDetectionURL = "http://localhost:8000/lips-movement-detection/"
	WhisperLargeV3URL        = "http://localhost:8000/whisper-large-v3/"
	OnePersonDetectURL       = "http://localhost:8000/one-person-detect/"
)

type processFunction func(client *http.Client, baseUrl string, video []byte, format string) (*domain.VideoFakeCandidat, error)

type ModelClientImpl struct {
	BaseURL         string
	processFunction processFunction
	HTTPClient      *http.Client
}

func NewModelClientImpl(baseURL string, processFunction processFunction) *ModelClientImpl {
	client := &ModelClientImpl{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 600 * time.Second,
		},
		processFunction: processFunction,
	}

	return client
}

func (c *ModelClientImpl) ping() error {
	resp, err := c.HTTPClient.Get(c.BaseURL) // Используйте конкретный endpoint для проверки доступности
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server: %s not ready, status code: %d", c.BaseURL, resp.StatusCode)
	}
	return nil
}

func (c *ModelClientImpl) Process(video []byte, format string) (*domain.VideoFakeCandidat, error) {
	return c.processFunction(c.HTTPClient, c.BaseURL, video, format)
}

// Создание клиентов для каждого сервиса.
var (
	DeepfakeDetectClient        = ModelClientImpl{BaseURL: DeepfakeDetectURL}
	OpenClosedEyeDetectClient   = ModelClientImpl{BaseURL: OpenClosedEyeDetectURL}
	AudioFakeDetectionClient    = ModelClientImpl{BaseURL: AudioFakeDetectionURL}
	LipsMovementDetectionClient = ModelClientImpl{BaseURL: LipsMovementDetectionURL}
	WhisperLargeV3Client        = ModelClientImpl{BaseURL: WhisperLargeV3URL}
	OnePersonDetectClient       = ModelClientImpl{BaseURL: OnePersonDetectURL}
)
