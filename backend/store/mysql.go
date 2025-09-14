package store

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Gopher0727/GoRepo/backend/models"
)

var DB *gorm.DB

// InitMySQL 初始化 GORM MySQL 连接
func InitMySQL(cfg models.MySQLConfig) (*gorm.DB, error) {
	user := cfg.User
	if user == "" {
		user = "root"
	}

	host := cfg.Host
	if host == "" {
		host = "127.0.0.1"
	}

	if cfg.Port == 0 {
		cfg.Port = 3306
	}

	if cfg.Database == "" {
		cfg.Database = "testdb"
	}

	{
		rawDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", user, cfg.RootPassword, host, cfg.Port)
		tmp, err := gorm.Open(mysql.Open(rawDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
		if err != nil {
			return nil, err
		}
		sqlDB, err := tmp.DB()
		if err != nil {
			return nil, err
		}
		if _, err := sqlDB.Exec("CREATE DATABASE IF NOT EXISTS `" + cfg.Database + "` DEFAULT CHARACTER SET utf8mb4"); err != nil {
			return nil, err
		}
		_ = sqlDB.Close()
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, cfg.RootPassword, host, cfg.Port, cfg.Database)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = db

	return db, nil
}
