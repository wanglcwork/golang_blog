package service

import (
	"blog-backend/model"
	"errors"

	"gorm.io/gorm"
	"github.com/sirupsen/logrus"
)

// PostService 文章服务接口
type PostService interface {
	CreatePost(title, content string, userID uint) (*model.Post, error)
	GetPostByID(id uint) (*model.Post, error)
	ListPosts(page, pageSize int) ([]model.Post, int64, error)
	UpdatePost(id uint, title, content string, userID uint) (*model.Post, error)
	DeletePost(id uint, userID uint) error
	GetUserPosts(userID uint, page, pageSize int) ([]model.Post, int64, error)
}

// postService 文章服务实现
type postService struct {
	db *gorm.DB
}

// NewPostService 创建文章服务实例
func NewPostService(db *gorm.DB) PostService {
	return &postService{db: db}
}

// CreatePost 创建文章
func (s *postService) CreatePost(title, content string, userID uint) (*model.Post, error) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := s.db.Create(post).Error; err != nil {
		logrus.Errorf("创建文章失败: %v", err)
		return nil, err
	}

	logrus.Infof("用户 %d 创建文章成功: %s", userID, title)
	return post, nil
}

// GetPostByID 根据ID获取文章
func (s *postService) GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := s.db.Preload("User").Preload("Comments").Preload("Comments.User").First(&post, id).Error; err != nil {
		logrus.Errorf("获取文章 %d 失败: %v", id, err)
		return nil, err
	}
	return &post, nil
}

// ListPosts 获取文章列表（分页）
func (s *postService) ListPosts(page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	// 计算总记录数
	if err := s.db.Model(&model.Post{}).Count(&total).Error; err != nil {
		logrus.Errorf("计算文章总数失败: %v", err)
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	if err := s.db.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts).Error; err != nil {
		logrus.Errorf("获取文章列表失败: %v", err)
		return nil, 0, err
	}

	return posts, total, nil
}

// UpdatePost 更新文章
func (s *postService) UpdatePost(id uint, title, content string, userID uint) (*model.Post, error) {
	// 检查文章是否存在
	var post model.Post
	if err := s.db.First(&post, id).Error; err != nil {
		logrus.Errorf("更新文章 %d 失败: 文章不存在 - %v", id, err)
		return nil, errors.New("文章不存在")
	}

	// 检查权限
	if post.UserID != userID {
		logrus.Warnf("用户 %d 尝试更新不属于自己的文章 %d", userID, id)
		return nil, errors.New("没有权限更新此文章")
	}

	// 更新文章
	post.Title = title
	post.Content = content
	if err := s.db.Save(&post).Error; err != nil {
		logrus.Errorf("更新文章 %d 失败: %v", id, err)
		return nil, err
	}

	logrus.Infof("用户 %d 更新文章成功: %d", userID, id)
	return &post, nil
}

// DeletePost 删除文章
func (s *postService) DeletePost(id uint, userID uint) error {
	// 检查文章是否存在
	var post model.Post
	if err := s.db.First(&post, id).Error; err != nil {
		logrus.Errorf("删除文章 %d 失败: 文章不存在 - %v", id, err)
		return errors.New("文章不存在")
	}

	// 检查权限
	if post.UserID != userID {
		logrus.Warnf("用户 %d 尝试删除不属于自己的文章 %d", userID, id)
		return errors.New("没有权限删除此文章")
	}

	// 删除文章（软删除）
	if err := s.db.Delete(&post).Error; err != nil {
		logrus.Errorf("删除文章 %d 失败: %v", id, err)
		return err
	}

	logrus.Infof("用户 %d 删除文章成功: %d", userID, id)
	return nil
}

// GetUserPosts 获取用户的文章列表
func (s *postService) GetUserPosts(userID uint, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	// 计算总记录数
	if err := s.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		logrus.Errorf("计算用户 %d 的文章总数失败: %v", userID, err)
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取分页数据
	if err := s.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts).Error; err != nil {
		logrus.Errorf("获取用户 %d 的文章列表失败: %v", userID, err)
		return nil, 0, err
	}

	return posts, total, nil
}
