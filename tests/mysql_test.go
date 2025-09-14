package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql" // 只执行包的 init()
	"github.com/spf13/viper"

	"github.com/Gopher0727/GoRepo/backend/models"
	"github.com/Gopher0727/GoRepo/backend/pkg/logger"
)

func TestMySQL(t *testing.T) {
	// 读取配置
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		logger.Error.Println(err)
		t.Fail()
		return
	}

	var cfg models.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error.Println(err)
		t.Fail()
		return
	}

	// 连接 MySQL
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/%s", cfg.MySQL.RootPassword, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error.Println(err)
		t.Fail()
		return
	}
	defer db.Close()

	// 创建测试表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL
	)`)
	if err != nil {
		logger.Error.Println(err)
		t.Fail()
		return
	}

	// 插入测试数据
	_, err = db.Exec(`INSERT INTO users (name) VALUES (?)`, "TestUser")
	if err != nil {
		logger.Error.Println(err)
		t.Fail()
		return
	}

	// 查询验证
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE name = ?", "TestUser").Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error.Println("测试失败：未找到用户记录")
		} else {
			logger.Error.Println(err)
		}
		t.Fail()
		return
	}

	logger.Info.Println("MySQL Test Passed, User:", name)

	// 清理测试数据
	_, _ = db.Exec("DELETE FROM users WHERE name = ?", "TestUser")
}
