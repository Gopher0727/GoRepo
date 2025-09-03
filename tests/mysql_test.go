package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql" // 只执行包的 init()
	"github.com/pelletier/go-toml/v2"

	"github.com/Gopher0727/GoRepo/backend/config"
)

func TestMySQL(t *testing.T) {
	// 读取配置
	data, err := os.ReadFile("../config.toml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg config.Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:%d)/%s", cfg.MySQL.RootPassword, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建测试表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// 插入测试数据
	_, err = db.Exec(`INSERT INTO users (name) VALUES (?)`, "TestUser")
	if err != nil {
		log.Fatal(err)
	}

	// 查询验证
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE name = ?", "TestUser").Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("测试失败：未找到用户记录")
			return
		}
		log.Fatal(err)
	}

	fmt.Println("MySQL Test Passed, User:", name)

	// 清理测试数据
	_, _ = db.Exec("DELETE FROM users WHERE name = ?", "TestUser")
}
