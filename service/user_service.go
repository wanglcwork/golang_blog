package service

import (
	"blog-backend/model"
	// "blog-backend/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/sirupsen/logrus"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(username, email, password string) (*model.User, error)
	Login(username, password string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

// CreateUser 创建新用户
func (s *userService) CreateUser(username, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		logrus.Warnf("用户名已存在: %s", username)
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		logrus.Warnf("邮箱已存在: %s", email)
		return nil, errors.New("邮箱已存在")
	}

	// 创建新用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: password, // 密码会在BeforeSave钩子中加密
	}

	if err := s.db.Create(user).Error; err != nil {
		logrus.Errorf("创建用户失败: %v", err)
		return nil, err
	}

	logrus.Infof("用户创建成功: %s", username)
	return user, nil
}

// Login 用户登录
func (s *userService) Login(username, password string) (*model.User, error) {
	// 查询用户
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		logrus.Warnf("用户不存在: %s", username)
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Warnf("用户 %s 密码错误", username)
		return nil, errors.New("用户名或密码错误")
	}

	logrus.Infof("用户登录成功: %s", username)
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		logrus.Errorf("获取用户 %d 失败: %v", id, err)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		logrus.Errorf("获取用户 %s 失败: %v", username, err)
		return nil, err
	}
	return &user, nil
}
