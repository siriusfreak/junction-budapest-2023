package tasksstorage

import (
	"context"
	"encoding/json"
	"orchestrator/internal/domain"
	"time"
)

func (s *Storage) AddOrUpdateTask(ctx context.Context, video *domain.VideoFakeCandidat) error {
	videoData, err := json.Marshal(video)
	if err != nil {
		return err
	}

	key := "task:" + video.UID

	exists, err := s.client.Exists(key).Result()
	if err != nil {
		return err
	}

	if exists > 0 {
		currentData, err := s.client.Get(key).Result()
		if err != nil {
			return err
		}

		var currentVideo domain.VideoFakeCandidat
		err = json.Unmarshal([]byte(currentData), &currentVideo)
		if err != nil {
			return err
		}

		if video.AudioFakeDetectionResult != nil {
			currentVideo.AudioFakeDetectionResult = video.AudioFakeDetectionResult
		}
		if video.DeepfakeDetectResult != nil {
			currentVideo.DeepfakeDetectResult = video.DeepfakeDetectResult
		}
		if video.LipsMovementDetectionResult != nil {
			currentVideo.LipsMovementDetectionResult = video.LipsMovementDetectionResult
		}
		if video.OpenClosedEyeDetect != nil {
			currentVideo.AudioFakeDetectionResult = video.OpenClosedEyeDetect
		}
		if video.WhisperLargeV3Result != nil {
			currentVideo.WhisperLargeV3Result = video.WhisperLargeV3Result
		}
		if video.OnePersonDetectResult != nil {
			currentVideo.OnePersonDetectResult = video.OnePersonDetectResult
		}

		videoData, err = json.Marshal(currentVideo)
		if err != nil {
			return err
		}
	}

	err = s.client.Set(key, videoData, time.Hour*12).Err()
	if err != nil {
		return err
	}

	return nil
}
