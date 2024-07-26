package model

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, data any, exp time.Duration) error
	Get(ctx context.Context, key string, data any) error
	Del(ctx context.Context, keys ...string) error
}
