package config

type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	ContainerName string `mapstructure:"container_name"`
	RootPassword  string `mapstructure:"root_password"`
	Database      string `mapstructure:"database"`
	Port          int    `mapstructure:"port"`
	Volume        string `mapstructure:"volume"`
}

type RedisConfig struct {
	ContainerName string `mapstructure:"container_name"`
	Password      string `mapstructure:"password"`
	Port          int    `mapstructure:"port"`
	Volume        string `mapstructure:"volume"`
}
