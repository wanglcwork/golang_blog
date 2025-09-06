package controller

import (
	"blog-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CommentController 评论控制器
type CommentController struct {
	commentService service.CommentService
}

// NewCommentController 创建评论控制器实例
func NewCommentController(commentService service.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		logrus.Warn("创建评论时未获取到用户ID")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取文章ID
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的文章ID: %s", postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required,min=1,max=500"`
	}

	// 绑定并验证输入
	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Warnf("创建评论输入验证失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的输入数据: " + err.Error()})
		return
	}

	// 创建评论
	comment, err := c.commentService.CreateComment(input.Content, userID.(uint), uint(postID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "文章不存在" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"user_id":    comment.UserID,
			"post_id":    comment.PostID,
			"created_at": comment.CreatedAt,
		},
	})
}

// GetPostComments 获取文章的评论列表
func (c *CommentController) GetPostComments(ctx *gin.Context) {
	// 获取文章ID
	postIDStr := ctx.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的文章ID: %s", postIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取分页参数
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取评论列表
	comments, total, err := c.commentService.GetPostComments(uint(postID), page, pageSize)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "文章不存在" {
			statusCode = http.StatusNotFound
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 处理评论数据
	var commentList []gin.H
	for _, comment := range comments {
		commentList = append(commentList, gin.H{
			"id":         comment.ID,
			"content":    comment.Content,
			"user_id":    comment.UserID,
			"username":   comment.User.Username,
			"created_at": comment.CreatedAt,
		})
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"comments":   commentList,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// DeleteComment 删除评论
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		logrus.Warn("删除评论时未获取到用户ID")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 获取评论ID
	commentIDStr := ctx.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		logrus.Warnf("无效的评论ID: %s", commentIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	// 删除评论
	err = c.commentService.DeleteComment(uint(commentID), userID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "评论不存在" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "没有权限删除此评论" {
			statusCode = http.StatusForbidden
		}
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}
