package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"

	"github.com/Gopher0727/GoRepo/backend/config"
)

func main() {
	// 读取 TOML 配置（根目录执行）
	data, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg config.Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	// 生成 docker-compose 文件，文件权限 0644（八进制数）
	composeContent := generateDockerCompose(cfg)
	if err := os.WriteFile("docker-compose.yml", []byte(composeContent), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("docker-compose.yml 已生成成功！")
}

func generateDockerCompose(cfg config.Config) string {
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
