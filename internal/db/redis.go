package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/labstack/gommon/log"
	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	redisClient *redis.Client
}

func NewRedisClient() *RedisClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &RedisClient{redisClient: redisClient}
}

func (r *RedisClient) Set(ctx context.Context, key string, value []byte, exp time.Duration) error {
	err := r.redisClient.Set(ctx, key, value, exp).Err()
	if err != nil {
		log.Errorf("failed insert data to redis, error: %v", err)
	}
	return err
}

func (r *RedisClient) Get(ctx context.Context, key string, data any) error {
	value, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Errorf("failed get data from redis, error: %v", err)
		return err
	}
	return json.Unmarshal([]byte(value), &data)
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	err := r.redisClient.Del(ctx, keys...).Err()
	if err != nil {
		log.Errorf("failed delete data from redis, error: %v", err)
	}
	return err
}
