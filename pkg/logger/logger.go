package logger

import (
	"io"
	"log"
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Trace   *log.Logger // 几乎任何东西
	Info    *log.Logger // 重要信息
	Warning *log.Logger // 警告信息
	Error   *log.Logger // 错误信息
)

// ANSI 颜色码
const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorYellow = "\x1b[33m"
	colorGreen  = "\x1b[32m"
)

func init() {
	// 创建日志文件
	fileName := "../logs/app.log"

	// file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("无法打开 %s 文件: %s", fileName, err)
	// }

	// * 使用 lumberjack 做文件轮转（按大小切割）
	file := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,    // 每个日志文件的最大大小（MB）
		MaxBackups: 7,     // 保留的旧日志文件数量
		MaxAge:     30,    // 保留日志的天数
		Compress:   false, // 是否压缩旧日志文件
	}

	Trace = log.New(
		&colorWriter{
			file:    file,
			console: io.Discard,
		},
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Info = log.New(
		&colorWriter{
			console: os.Stdout,
			color:   colorGreen,
			reset:   colorReset,
		},
		"[INFO]: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Warning = log.New(
		&colorWriter{
			console: os.Stdout,
			color:   colorYellow,
			reset:   colorReset,
		},
		"[WARNING]: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Error = log.New(
		// 将日志同时写入文件和标准错误输出
		&colorWriter{
			file:    file,
			console: os.Stderr,
			color:   colorRed,
			reset:   colorReset,
		},
		"[ERROR]: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}

// 标准库的 log.Logger 对单个 Logger 的写操作有互斥保护，
// 但包中有多个 Logger（Info、Warning、Error）共享同一个 colorWriter，
// 多个不同 Logger 同时写到同一 colorWriter 时可能发生行间交错。

type colorWriter struct {
	mu      sync.Mutex // 确保线程安全
	file    io.Writer
	console io.Writer
	color   string
	reset   string
}

func (w *colorWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 写入文件（无颜色）
	if w.file != nil {
		if _, err := w.file.Write(p); err != nil {
			return 0, err
		}
	}
	// 写入终端（带颜色包裹）
	if w.console != nil {
		colored := append([]byte(w.color), p...)
		colored = append(colored, []byte(w.reset)...)
		if _, err := w.console.Write(colored); err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
