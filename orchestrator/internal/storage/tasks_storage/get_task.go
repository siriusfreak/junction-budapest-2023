package tasksstorage

import (
	"context"
	"encoding/json"
	"orchestrator/internal/domain"

	"github.com/go-redis/redis"
)

func (s *Storage) GetTask(ctx context.Context, uid string) (*domain.VideoFakeCandidat, error) {
	key := "task:" + uid

	videoData, err := s.client.Get(key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var video domain.VideoFakeCandidat
	err = json.Unmarshal([]byte(videoData), &video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}
