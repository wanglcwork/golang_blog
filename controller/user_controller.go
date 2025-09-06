package controller

import (
	"blog-backend/config"
	"blog-backend/service"
	"blog-backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
	cfg         *config.Config
}

// NewUserController 创建用户控制器实例
func NewUserController(userService service.UserService, cfg *config.Config) *UserController {
	return &UserController{
		userService: userService,
		cfg:         cfg,
	}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// 绑定并验证输入
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("注册输入验证失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据: " + err.Error()})
		return
	}

	// 创建用户
	user, err := c.userService.CreateUser(input.Username, input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "用户注册成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定并验证输入
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("登录输入验证失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据: " + err.Error()})
		return
	}

	// 验证用户
	user, err := c.userService.Login(input.Username, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, c.cfg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// GetUser 获取用户信息
func (c *UserController) GetUser(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的用户ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取用户信息
	user, err := c.userService.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
