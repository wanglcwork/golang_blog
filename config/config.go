package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config 应用配置
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBCharset  string
	JWTSecret  string
	JWTExpiry  string
	ServerPort string
	GinMode    string
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("未找到 .env 文件，将使用环境变量")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBName:     getEnv("DB_NAME", "blog_db"),
		DBCharset:  getEnv("DB_CHARSET", "utf8mb4"),
		JWTSecret:  getEnv("JWT_SECRET", "default_secret"),
		JWTExpiry:  getEnv("JWT_EXPIRATION", "720h"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}

	return config, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
