package middleware

import (
	"blog-backend/config"
	"blog-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("请求头中Authorization为空")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logrus.Warn("Authorization格式错误")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			c.Abort()
			return
		}

		// 解析JWT
		tokenString := parts[1]
		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			logrus.Warnf("JWT解析错误: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
