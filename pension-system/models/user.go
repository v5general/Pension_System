package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"not null" binding:"required"` // 用户名，可以重复
	Password  string         `json:"-" gorm:"not null" binding:"required"`
	Name      string         `json:"name" gorm:"not null;default:''" binding:"required"` // 账户名，与活跃用户不重复（已删除用户的账户名可重用）
	Role      string         `json:"role" gorm:"not null;default:'user'"` // admin, operator, user
	LastSeen  time.Time      `json:"last_seen"` // Last activity time for online status
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
