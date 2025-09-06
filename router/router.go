package router

import (
	"blog-backend/config"
	"blog-backend/controller"
	"blog-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	userController *controller.UserController,
	postController *controller.PostController,
	commentController *controller.CommentController,
	cfg *config.Config,
) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	// API路由组
	api := r.Group("/api")
	{
		// 公共路由
		public := api.Group("")
		{
			// 用户相关
			public.POST("/register", userController.Register)
			public.POST("/login", userController.Login)
			public.GET("/users/:id", userController.GetUser)

			// 文章相关
			public.GET("/posts", postController.ListPosts)
			public.GET("/posts/:id", postController.GetPost)
			public.GET("/users-posts/:user_id/posts", postController.GetUserPosts)

			// 评论相关
			public.GET("/posts-comments/:post_id/comments", commentController.GetPostComments)
		}

		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// 文章相关
			protected.POST("/posts", postController.CreatePost)
			protected.PUT("/posts/:id", postController.UpdatePost)
			protected.DELETE("/posts/:id", postController.DeletePost)

			// 评论相关
			protected.POST("/posts-comments/:post_id/comments", commentController.CreateComment)
			protected.DELETE("/comments/:id", commentController.DeleteComment)
		}
	}

	return r
}
