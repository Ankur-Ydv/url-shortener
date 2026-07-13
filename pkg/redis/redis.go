package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(context context.Context, redisHost string, redisPort string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       redisHost + ":" + redisPort,
		MaxRetries: 3,
	})

	if err := client.Ping(context).Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return client, nil
}

func CloseRedisClient(context context.Context, client *redis.Client) {
	_ = client.Close()
}

func SetKey(context context.Context, client *redis.Client, key string, value string) error {
	if err := client.Set(context, key, value, 0).Err(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func GetKey(context context.Context, client *redis.Client, key string) (string, error) {
	value, err := client.Get(context, key).Result()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return value, nil
}

func DeleteKey(context context.Context, client *redis.Client, key string) error {
	if err := client.Del(context, key).Err(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
