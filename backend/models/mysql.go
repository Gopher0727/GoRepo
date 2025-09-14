package models

// MySQLConfig 数据库连接配置
// Host 为空默认 127.0.0.1，User 为空默认 root
// Port 为空默认 3306
// RootPassword 也作为普通账号密码使用（若未拆分用户）
type MySQLConfig struct {
	ContainerName string `mapstructure:"container_name"`
	RootPassword  string `mapstructure:"root_password"`
	Database      string `mapstructure:"database"`
	Port          int    `mapstructure:"port"`
	Volume        string `mapstructure:"volume"`
	Host          string `mapstructure:"host"`
	User          string `mapstructure:"user"`
}
