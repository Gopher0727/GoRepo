package main

import (
	"log"
	"testing"

	"github.com/Gopher0727/GoRepo/backend/pkg/logger"
)

func TestLogger(t *testing.T) {
	log.Println("日志系统已初始化，日志将写入 app.log 文件")
	// 示例日志
	logger.Trace.Println("这是一个跟踪信息")
	logger.Info.Println("这是一个普通信息")
	logger.Warning.Println("这是一个警告信息")
	logger.Error.Println("这是一个错误信息")
}
