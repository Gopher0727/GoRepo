package config

type Config struct {
	MySQL MySQLConfig `toml:"mysql"`
	Redis RedisConfig `toml:"redis"`
}

type MySQLConfig struct {
	ContainerName string `toml:"container_name"`
	RootPassword  string `toml:"root_password"`
	Database      string `toml:"database"`
	Port          int    `toml:"port"`
	Volume        string `toml:"volume"`
}

type RedisConfig struct {
	ContainerName string `toml:"container_name"`
	Password      string `toml:"password"`
	Port          int    `toml:"port"`
	Volume        string `toml:"volume"`
}
