package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null" binding:"required"`
	Password  string         `json:"-" gorm:"not null" binding:"required"`
	Name      string         `json:"name" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;default:'user'"` // admin, operator, user
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
