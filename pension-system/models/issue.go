package models

import (
	"time"

	"gorm.io/gorm"
)

type Issue struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null" binding:"required"`
	Description string         `json:"description"`
	Category    string         `json:"category" gorm:"not null"` // 生活服务, 医疗健康, 文化娱乐, 其他
	Priority    string         `json:"priority" gorm:"default:'normal'"` // low, normal, high, urgent
	Status      string         `json:"status" gorm:"default:'pending'"` // pending, processing, resolved, closed
	SubmitterID uint           `json:"submitter_id"`
	Submitter   *User          `json:"submitter,omitempty" gorm:"foreignKey:SubmitterID"`
	HandlerID   *uint          `json:"handler_id"`
	Handler     *User          `json:"handler,omitempty" gorm:"foreignKey:HandlerID"`
	Response    string         `json:"response"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Issue) TableName() string {
	return "issues"
}
