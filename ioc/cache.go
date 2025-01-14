package ioc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	conf "muxi_auditor/config"
)

func InitCache(cfg *conf.CacheConfig) *redis.Client {
	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // Redis 地址
		Password: cfg.Password, // Redis 密码
	})

	// 测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %w", err))
	}

	return client
}
