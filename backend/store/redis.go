package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Gopher0727/GoRepo/backend/models"
)

// RDB 全局 Redis 客户端
var RDB *redis.Client

// InitRedis 初始化 Redis 连接
// 根据 models.RedisConfig 设置默认值并进行健康检查
func InitRedis(cfg models.RedisConfig) error {
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}

	if cfg.Port == 0 {
		cfg.Port = 6379
	}

	// DB 默认 0（结构体标签注释里已说明）

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}
	RDB = rdb

	return nil
}
