package models

// Config 应用整体配置
// 通过 viper.Unmarshal(&cfg) 填充
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
	HTTP  HTTPConfig  `mapstructure:"http"`
}

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

// HTTPConfig HTTP 服务配置
// Port 默认 8080; JWTSecret 可在开发缺省
// 在生产环境务必通过环境变量或安全方式覆盖
// (例如: export APP_HTTP__JWT_SECRET=xxxx 并用 viper.AutomaticEnv)
type HTTPConfig struct {
	Port      int    `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
}
