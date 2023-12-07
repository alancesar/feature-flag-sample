package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type (
	Redis struct {
		client *redis.Client
	}
)

func NewRedis(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func (r Redis) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := r.client.Get(ctx, key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return cmd.Bytes()
}

func (r Redis) Set(ctx context.Context, key string, value []byte) error {
	return r.client.Set(ctx, key, value, time.Minute).Err()
}
