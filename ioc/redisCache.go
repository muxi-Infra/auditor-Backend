package ioc

import (
	"context"
	"encoding/json"
	errrs "errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/muxi-Infra/auditor-Backend/repository/cache/errorxs"
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

// GetStringSlice 获取tag
func (rc *RedisCache) GetStringSlice(ctx context.Context, key string) ([]string, error) {
	val, err := rc.Client.Get(ctx, key).Result()
	if err != nil {
		if errrs.Is(err, redis.Nil) {
			return nil, errorxs.ToCacheNotFoundError(err)
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var result []string
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return result, nil
}

// SetStringSlice 设置tag
func (rc *RedisCache) SetStringSlice(ctx context.Context, key string, val []string, expiration time.Duration) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	return rc.Client.Set(ctx, key, bytes, expiration).Err()
}

func (rc *RedisCache) SetString(ctx context.Context, key string, val string, expiration time.Duration) error {
	return rc.Client.Set(ctx, key, val, expiration).Err()
}

func (rc *RedisCache) GetString(ctx context.Context, key string) (string, error) {
	val, err := rc.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
