package main

import (
	"blog-backend/config"
	"blog-backend/controller"
	"blog-backend/model"
	"blog-backend/router"
	"blog-backend/service"
	"blog-backend/utils"
	"fmt"

	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("无法加载配置: %v", err)
	}

	// 初始化日志
	utils.InitLogger()

	// 连接数据库
	db, err := model.InitDB(cfg)
	if err != nil {
		logrus.Fatalf("数据库连接失败: %v", err)
	}
	defer model.CloseDB(db)

	// 自动迁移数据表
	model.AutoMigrate(db)

	// 初始化服务
	userService := service.NewUserService(db)
	postService := service.NewPostService(db)
	commentService := service.NewCommentService(db)

	// 初始化控制器
	userController := controller.NewUserController(userService, cfg)
	postController := controller.NewPostController(postService)
	commentController := controller.NewCommentController(commentService)

	// 设置路由
	r := router.SetupRouter(userController, postController, commentController, cfg)

	// 启动服务器
	logrus.Printf("服务器启动在端口 %s", cfg.ServerPort)
	logrus.Fatal(r.Run(fmt.Sprintf(":%s", cfg.ServerPort)))
}
