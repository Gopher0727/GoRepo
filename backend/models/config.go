package models

// Config 应用配置
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	HTTP  HTTPConfig  `mapstructure:"http"`
}
