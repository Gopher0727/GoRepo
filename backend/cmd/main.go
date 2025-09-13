package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/Gopher0727/GoRepo/backend/config"
)

func main() {
	// 读取 TOML 配置（根目录执行）

	// config 文件中用 mapstructure 指定
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", err))
	}

	// config 文件中用 toml 指定
	// data, err := os.ReadFile("config.toml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// var cfg config.Config
	// if err := toml.Unmarshal(data, &cfg); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println(cfg.MySQL.ContainerName)
	fmt.Println(cfg.MySQL.RootPassword)
	fmt.Println(cfg.MySQL.Database)
	fmt.Println(cfg.MySQL.Port)
	fmt.Println(cfg.MySQL.Volume)

	fmt.Println(cfg.Redis.ContainerName)
	fmt.Println(cfg.Redis.Password)
	fmt.Println(cfg.Redis.Port)
	fmt.Println(cfg.Redis.Volume)

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
