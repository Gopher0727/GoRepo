package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/Gopher0727/GoRepo/backend/models"
)

func main() {
	cfg := loadConfig()
	content := generateDockerCompose(cfg)
	if err := os.WriteFile("docker-compose.yml", []byte(content), 0644); err != nil {
		log.Fatalf("写入 docker-compose.yml 失败: %v", err)
	}
	log.Println("docker-compose.yml 已生成")
}

func loadConfig() models.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	var cfg models.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	return cfg
}

func generateDockerCompose(cfg models.Config) string {
	return fmt.Sprintf(`services:
  mysql:
    image: mysql:8.0
    container_name: %s
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: %s
      MYSQL_DATABASE: %s
    ports:
      - "%d:3306"
    volumes:
      - %s

  redis:
    image: redis:7.0
    container_name: %s
    restart: unless-stopped
    command: ["redis-server", "--requirepass", "%s"]
    ports:
      - "%d:6379"
    volumes:
      - %s

volumes:
  mysql_data:
  redis_data:
`, cfg.MySQL.ContainerName, cfg.MySQL.RootPassword, cfg.MySQL.Database, cfg.MySQL.Port, cfg.MySQL.Volume,
		cfg.Redis.ContainerName, cfg.Redis.Password, cfg.Redis.Port, cfg.Redis.Volume)
}
