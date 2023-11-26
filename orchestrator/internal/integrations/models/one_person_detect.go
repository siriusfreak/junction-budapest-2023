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

func OnePersonDetectProcess(client *http.Client, baseUrl string, video []byte) (*domain.VideoFakeCandidat, error) {
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

	_ = writer.WriteField("processed_percent", fmt.Sprintf("%d", 50))
	_ = writer.WriteField("confidence_threshold", fmt.Sprintf("%.2f", 0.3))
	_ = writer.WriteField("skip_milliseconds", fmt.Sprintf("%d", 1000))

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
		return nil, fmt.Errorf("server returned non-200 status code from OnePersonDetect: %d", resp.StatusCode)
	}

	var result struct {
		Frames          []string `json:"frames"`
		TotalFrames     int      `json:"total_frames"`
		ProcessedFrames int      `json:"processed_frames"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response from OnePersonDetect: %w", err)
	}

	onePersonDetected := len(result.Frames) == 0

	return &domain.VideoFakeCandidat{
		OnePersonDetectResult: &onePersonDetected,
	}, nil
}
