package models

import (
	"fmt"
	"log"
	"net/http"
	"orchestrator/internal/domain"
	"time"
)

// Определение базового URL для каждого сервиса.
const (
	DeepfakeDetectURL        = "http://deepfake-detect:8000"
	OpenClosedEyeDetectURL   = "http://open-closed-eye-detect:8000"
	AudioFakeDetectionURL    = "http://audio-fake-detection:8000"
	LipsMovementDetectionURL = "http://lips-movement-detection:8000"
	WhisperLargeV3URL        = "http://whisper-large-v3:8000"
	OnePersonDetectURL       = "http://one-person-detect:8000"
)

type processFunction func(client *http.Client, baseUrl string, video []byte) (*domain.VideoFakeCandidat, error)

type ModelClientImpl struct {
	BaseURL         string
	processFunction processFunction
	HTTPClient      *http.Client
}

func NewModelClientImpl(baseURL string, processFunction processFunction) *ModelClientImpl {
	client := &ModelClientImpl{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		processFunction: processFunction,
	}

	err := client.ping()
	if err != nil {
		log.Fatalf("cannot reach the server: %v", err)
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

func (c *ModelClientImpl) Process(video []byte) (*domain.VideoFakeCandidat, error) {
	return c.processFunction(c.HTTPClient, c.BaseURL, video)
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
