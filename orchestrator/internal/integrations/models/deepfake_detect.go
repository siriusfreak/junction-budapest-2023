package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"orchestrator/internal/domain"
)

func DeepfakeDetectProcess(client *http.Client, baseUrl string, video []byte) (*domain.VideoFakeCandidat, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("video", "video.mp4")
	if err != nil {
		return nil, err
	}
	_, err = part.Write(video)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", baseUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code from DeepfakeDetect: %d", resp.StatusCode)
	}

	var result struct {
		Fake float64 `json:"fake"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response from OnePersonDetect: %w", err)
	}

	deepfakeDetectResult := result.Fake <= 0.13
	return &domain.VideoFakeCandidat{
		DeepfakeDetectResult: &deepfakeDetectResult,
	}, nil
}
