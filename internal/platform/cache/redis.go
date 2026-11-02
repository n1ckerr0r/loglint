package cache

import (
	"context"
	"time"

	"github.com/n1ckerr0r/loglint/internal/service"
	"github.com/redis/go-redis/v9"
)

type NoopStatusCache struct{}

func (NoopStatusCache) SetStatus(context.Context, string, service.RunStatus) error {
	return nil
}

func (NoopStatusCache) Close() error {
	return nil
}

type RedisStatusCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStatusCache(redisURL string) (*RedisStatusCache, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &RedisStatusCache{
		client: redis.NewClient(options),
		ttl:    24 * time.Hour,
	}, nil
}

func (cache *RedisStatusCache) SetStatus(ctx context.Context, runID string, status service.RunStatus) error {
	return cache.client.Set(ctx, "loglint:run:"+runID+":status", string(status), cache.ttl).Err()
}

func (cache *RedisStatusCache) Close() error {
	return cache.client.Close()
}
