package service

import (
	"blog-backend/model"
	"errors"

	"gorm.io/gorm"
	"github.com/sirupsen/logrus"
)

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(content string, userID, postID uint) (*model.Comment, error)
	GetCommentByID(id uint) (*model.Comment, error)
	GetPostComments(postID uint, page, pageSize int) ([]model.Comment, int64, error)
	DeleteComment(id uint, userID uint) error
}

// commentService 评论服务实现
type commentService struct {
	db *gorm.DB
}

// NewCommentService 创建评论服务实例
func NewCommentService(db *gorm.DB) CommentService {
	return &commentService{db: db}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(content string, userID, postID uint) (*model.Comment, error) {
	// 检查文章是否存在
	var post model.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		logrus.Errorf("创建评论失败: 文章 %d 不存在 - %v", postID, err)
		return nil, errors.New("文章不存在")
	}

	// 创建评论
	comment := &model.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}

	if err := s.db.Create(comment).Error; err != nil {
		logrus.Errorf("创建评论失败: %v", err)
		return nil, err
	}

	logrus.Infof("用户 %d 为文章 %d 创建评论成功", userID, postID)
	return comment, nil
}

// GetCommentByID 根据ID获取评论
func (s *commentService) GetCommentByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := s.db.First(&comment, id).Error; err != nil {
		logrus.Errorf("获取评论 %d 失败: %v", id, err)
		return nil, err
	}
	return &comment, nil
}

// GetPostComments 获取文章的评论列表
func (s *commentService) GetPostComments(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 检查文章是否存在
	var post model.Post
	if err := s.db.First(&post, postID).Error; err != nil {
		logrus.Errorf("获取评论失败: 文章 %d 不存在 - %v", postID, err)
		return nil, 0, errors.New("文章不存在")
	}

	// 计算总记录数
	if err := s.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&total).Error; err != nil {
		logrus.Errorf("计算文章 %d 的评论总数失败: %v", postID, err)
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	if err := s.db.Preload("User").Where("post_id = ?", postID).Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error; err != nil {
		logrus.Errorf("获取文章 %d 的评论列表失败: %v", postID, err)
		return nil, 0, err
	}

	return comments, total, nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(id uint, userID uint) error {
	// 检查评论是否存在
	var comment model.Comment
	if err := s.db.First(&comment, id).Error; err != nil {
		logrus.Errorf("删除评论 %d 失败: 评论不存在 - %v", id, err)
		return errors.New("评论不存在")
	}

	// 检查权限
	if comment.UserID != userID {
		logrus.Warnf("用户 %d 尝试删除不属于自己的评论 %d", userID, id)
		return errors.New("没有权限删除此评论")
	}

	// 删除评论（软删除）
	if err := s.db.Delete(&comment).Error; err != nil {
		logrus.Errorf("删除评论 %d 失败: %v", id, err)
		return err
	}

	logrus.Infof("用户 %d 删除评论成功: %d", userID, id)
	return nil
}
