package model

import (
	"blog-backend/config"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接
var DB *gorm.DB

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"` // 不在JSON中显示密码
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Posts     []Post         `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Post 文章模型
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Comment 评论模型
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PostID    uint           `gorm:"not null" json:"post_id"`
	Post      Post           `gorm:"foreignKey:PostID" json:"post,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeSave - 保存前的钩子，用于密码加密
func (u *User) BeforeSave(tx *gorm.DB) error {
	if len(u.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=" + cfg.DBCharset + "&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Post{},
		&Comment{},
	)
}
