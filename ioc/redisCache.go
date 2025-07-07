package ioc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(c *redis.Client) *RedisCache {
	return &RedisCache{c}
}
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.Client.Set(ctx, key, value, expiration).Err()
}
func (rc *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return rc.Client.Get(ctx, key).Result()
}
