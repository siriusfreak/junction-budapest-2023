package models

import (
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

type processFunction func(client *HttpClientWithRetry, baseUrl string, video []byte, format string) (*domain.VideoFakeCandidat, error)

type HttpClientWithRetry struct {
	HTTPClient *http.Client
	MaxRetries int
	RetryDelay time.Duration
}

func (c *HttpClientWithRetry) Do(req *http.Request) (*http.Response, error) {
	var lastErr error
	for i := 0; i < c.MaxRetries; i++ {
		resp, err := c.HTTPClient.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(c.RetryDelay)
	}
	return nil, lastErr // return the last error encountered
}

type ModelClientImpl struct {
	BaseURL         string
	processFunction processFunction
	HTTPClient      *HttpClientWithRetry
}

func NewModelClientImpl(baseURL string, processFunction processFunction) *ModelClientImpl {
	client := &ModelClientImpl{
		BaseURL: baseURL,
		HTTPClient: &HttpClientWithRetry{
			HTTPClient: &http.Client{
				Timeout: 600 * time.Second,
			},
			MaxRetries: 5, // Define the maximum number of retries
			RetryDelay: 1 * time.Second, // Define the delay between retries
		},
		processFunction: processFunction,
	}

	return client
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
