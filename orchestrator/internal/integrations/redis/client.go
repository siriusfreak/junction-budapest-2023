package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

func CreateRedisClient(ctx context.Context, addr, pass string, number_db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       number_db,
	})
	if client == nil {
		log.Fatal("client is nil")
	}

	return client
}
