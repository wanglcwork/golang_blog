package utils

import (
	"blog-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// Claims JWT声明
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, cfg *config.Config) (string, error) {
	// 解析过期时间
	expirationTime, err := time.ParseDuration(cfg.JWTExpiry)
	if err != nil {
		logrus.Errorf("解析JWT过期时间错误: %v", err)
		return "", err
	}

	// 创建声明
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-backend",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		logrus.Errorf("生成JWT令牌错误: %v", err)
		return "", err
	}

	return tokenString, nil
}
