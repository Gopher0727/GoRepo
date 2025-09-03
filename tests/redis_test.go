package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/redis/go-redis/v9"

	"github.com/Gopher0727/GoRepo/backend/config"
)

func TestRedis(t *testing.T) {
	// 读取配置
	data, err := os.ReadFile("../config.toml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg config.Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%d", cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 写入测试数据
	err = rdb.Set(ctx, "test_key", "Hello, Redis", 10*time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}

	// 读取验证
	val, err := rdb.Get(ctx, "test_key").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Redis Test Failed: Key does not exist")
			return
		}
		log.Fatal(err)
	}

	fmt.Println("Redis Test Passed, test_key =", val)

	// 清理测试数据
	_ = rdb.Del(ctx, "test_key").Err()
}
