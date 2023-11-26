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
		log.Fatal("redis is nil")
	}

	res := client.Ping()
	if res.Err() != nil {
		log.Fatalf("redis ping err: %v", res.Err())
	}
	if _, err := res.Result(); err != nil {
		log.Fatalf("redis ping result err: %v", err)
	}

	return client
}
