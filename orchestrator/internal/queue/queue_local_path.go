package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Queue struct {
	name  string
	redis *redis.Client
}

func NewQueue(ctx context.Context, name string, redis *redis.Client) *Queue {
	return &Queue{
		name:  name,
		redis: redis,
	}
}

func (q *Queue) Pop(ctx context.Context) (string, error) {
	res, err := q.redis.BLPop(5*time.Second, q.name).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("error with Pop from queue: %s, err: %w", q.name, err)
	}

	return res[1], nil
}

func (q *Queue) Add(ctx context.Context, uid string) error {
	err := q.redis.LPush(q.name, uid).Err()
	if err != nil {
		return fmt.Errorf("error with Add in queue: %s, err: %w", q.name, err)
	}

	return nil
}
