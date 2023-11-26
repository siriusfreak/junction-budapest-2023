package models

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"orchestrator/internal/domain"
	"strings"
)

func WishperLargeV3Process(client *http.Client, baseUrl string, video []byte, format string) (*domain.VideoFakeCandidat, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("video", "video"+format)
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
		return nil, fmt.Errorf("server returned non-200 status code from WishperLargeV3: %d", resp.StatusCode)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	result := string(responseData)

	result = strings.ReplaceAll(result, ",", "")
	result = strings.ReplaceAll(result, " ", "")
	result = strings.ReplaceAll(result, ".", "")
	result = strings.ReplaceAll(result, "\"", "")
	result = strings.ToLower(result)
	result = strings.TrimSpace(result)

	whisperLargeV3Result := result == "twozerotwothree" || result == "2023"

	return &domain.VideoFakeCandidat{
		WhisperLargeV3Result: &whisperLargeV3Result,
	}, nil
}