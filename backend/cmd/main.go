package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/Gopher0727/GoRepo/backend/api"
	"github.com/Gopher0727/GoRepo/backend/middleware"
	"github.com/Gopher0727/GoRepo/backend/models"
	"github.com/Gopher0727/GoRepo/backend/pkg/logger"
	"github.com/Gopher0727/GoRepo/backend/store"
)

func loadConfig() models.Config {
	// 读取配置
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error.Printf("读取配置文件失败: %v", err)
		return models.Config{}
	}

	var cfg models.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error.Printf("解析配置文件失败: %v", err)
		return models.Config{}
	}

	if cfg.HTTP.Port == 0 {
		cfg.HTTP.Port = 8080
	}

	if cfg.HTTP.JWTSecret == "" {
		cfg.HTTP.JWTSecret = "dev-secret-change" // 开发缺省
	}

	return cfg
}

func main() {
	cfg := loadConfig()

	// 初始化 MySQL
	db, err := store.InitMySQL(cfg.MySQL)
	if err != nil {
		logger.Error.Printf("init mysql error: %v", err)
		return
	}

	// 自动迁移
	if err := db.AutoMigrate(&models.User{}, &models.Entry{}); err != nil {
		logger.Warning.Printf("auto migrate error: %v", err)
	}

	// 初始化 Redis
	if err := store.InitRedis(cfg.Redis); err != nil {
		logger.Error.Printf("init redis error: %v", err)
		return
	}

	// 初始化 JWT
	middleware.InitAuth(cfg.HTTP.JWTSecret)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	api.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: r,
	}

	go func() {
		logger.Info.Printf("HTTP server listening on :%d", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error.Fatalf("listen error: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info.Println("shutting down server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Warning.Printf("server shutdown error: %v", err)
	}
	logger.Info.Println("server exited")
}
