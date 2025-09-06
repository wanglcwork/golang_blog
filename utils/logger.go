package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// InitLogger 初始化日志配置
func InitLogger() {
	// 设置日志格式为JSON格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Printf("无法创建日志目录: %v", err)
		return
	}

	// 设置日志文件
	logFileName := time.Now().Format("2006-01-02") + ".log"
	logFilePath := filepath.Join(logDir, logFileName)
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Printf("无法打开日志文件: %v", err)
		return
	}

	// 设置输出为文件和控制台
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.InfoLevel)

	// 添加调用者信息
	logrus.SetReportCaller(true)
}
