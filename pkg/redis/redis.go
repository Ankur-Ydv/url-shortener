package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(context context.Context, redisHost string, redisPort string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	if err := client.Ping(context).Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return client, nil
}

func CloseRedisClient(context context.Context, client *redis.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
