package ioc

import (
	"context"
	"fmt"
	conf "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *conf.CacheConfig) *redis.Client {

	// 初始化 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // Redis 地址
		Password: cfg.Password, // Redis 密码，默认为空字符串
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		// 使用 fmt.Errorf 包装错误，便于追踪
		panic(fmt.Errorf("failed to connect to Redis: %v", err))
	}

	return client
}
