package models

// RedisConfig Redis 连接配置
// Host 为空默认 127.0.0.1，DB 默认为 0
type RedisConfig struct {
	ContainerName string `mapstructure:"container_name"`
	Password      string `mapstructure:"password"`
	Port          int    `mapstructure:"port"`
	Volume        string `mapstructure:"volume"`
	Host          string `mapstructure:"host"`
	DB            int    `mapstructure:"db"`
}
