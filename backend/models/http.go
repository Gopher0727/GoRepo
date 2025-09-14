package models

// HTTPConfig HTTP 服务配置
// Port 默认 8080; JWTSecret 可在开发缺省
// 在生产环境务必通过环境变量或安全方式覆盖
// (例如: export APP_HTTP__JWT_SECRET=xxxx 并用 viper.AutomaticEnv)
type HTTPConfig struct {
	Port      int    `mapstructure:"port"`
	JWTSecret string `mapstructure:"jwt_secret"`
}
