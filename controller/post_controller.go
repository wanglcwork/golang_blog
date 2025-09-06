package controller

import (
	"blog-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PostController 文章控制器
type PostController struct {
	postService service.PostService
}

// NewPostController 创建文章控制器实例
func NewPostController(postService service.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// CreatePost 创建文章
func (c *PostController) CreatePost(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		logrus.Warn("创建文章时未获取到用户ID")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required,min=3,max=100"`
		Content string `json:"content" binding:"required,min=10"`
	}

	// 绑定并验证输入
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("创建文章输入验证失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据: " + err.Error()})
		return
	}

	// 创建文章
	post, err := c.postService.CreatePost(input.Title, input.Content, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"user_id":    post.UserID,
			"created_at": post.CreatedAt,
		},
	})
}

// GetPost 获取单篇文章
func (c *PostController) GetPost(ctx *gin.Context) {
	// 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的文章ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取文章信息
	post, err := c.postService.GetPostByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 处理评论数据
	var comments []gin.H
	for _, comment := range post.Comments {
		comments = append(comments, gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"user_id":    comment.UserID,
			"username":   comment.User.Username,
			"created_at": comment.CreatedAt,
		})
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"user_id":    post.UserID,
		"username":   post.User.Username,
		"created_at": post.CreatedAt,
		"updated_at": post.UpdatedAt,
		"comments":   comments,
	})
}

// ListPosts 获取文章列表
func (c *PostController) ListPosts(ctx *gin.Context) {
	// 获取分页参数
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取文章列表
	posts, total, err := c.postService.ListPosts(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败: " + err.Error()})
		return
	}

	// 处理文章数据
	var postList []gin.H
	for _, post := range posts {
		postList = append(postList, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"user_id":    post.UserID,
			"username":   post.User.Username,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"posts":      postList,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// UpdatePost 更新文章
func (c *PostController) UpdatePost(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		logrus.Warn("更新文章时未获取到用户ID")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的文章ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required,min=3,max=100"`
		Content string `json:"content" binding:"required,min=10"`
	}

	// 绑定并验证输入
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("更新文章输入验证失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据: " + err.Error()})
		return
	}

	// 更新文章
	post, err := c.postService.UpdatePost(uint(id), input.Title, input.Content, userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "文章不存在" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "没有权限更新此文章" {
			statusCode = http.StatusForbidden
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post": gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"content":    post.Content,
			"updated_at": post.UpdatedAt,
		},
	})
}

// DeletePost 删除文章
func (c *PostController) DeletePost(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		logrus.Warn("删除文章时未获取到用户ID")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取文章ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的文章ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 删除文章
	err = c.postService.DeletePost(uint(id), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "文章不存在" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "没有权限删除此文章" {
			statusCode = http.StatusForbidden
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

// GetUserPosts 获取用户的文章列表
func (c *PostController) GetUserPosts(ctx *gin.Context) {
	// 获取用户ID
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的用户ID: %s", userIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取分页参数
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 获取用户的文章列表
	posts, total, err := c.postService.GetUserPosts(uint(userID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户文章列表失败: " + err.Error()})
		return
	}

	// 处理文章数据
	var postList []gin.H
	for _, post := range posts {
		postList = append(postList, gin.H{
			"id":         post.ID,
			"title":      post.Title,
			"created_at": post.CreatedAt,
			"updated_at": post.UpdatedAt,
		})
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"posts":      postList,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}
