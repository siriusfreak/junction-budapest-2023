package tasksstorage

import (
	"context"

	"github.com/go-redis/redis"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(ctx context.Context, client *redis.Client) *Storage {
	return &Storage{
		client: client,
	}
}
