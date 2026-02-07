package models

import (
	"time"

	"gorm.io/gorm"
)

type Survey struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null" binding:"required"`
	Description string         `json:"description"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Status      string         `json:"status" gorm:"default:'draft'"` // draft, active, closed
	Options     string         `json:"options" gorm:"type:text"` // JSON string of options
	TotalVotes  int            `json:"total_votes" gorm:"default:0"`
	CreatedBy   uint           `json:"created_by"`
	Creator     *User          `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Survey) TableName() string {
	return "surveys"
}

type SurveyVote struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	SurveyID  uint           `json:"survey_id" gorm:"not null"`
	Survey    *Survey        `json:"survey,omitempty" gorm:"foreignKey:SurveyID"`
	Option    string         `json:"option" gorm:"not null"`
	UserID    uint           `json:"user_id"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (SurveyVote) TableName() string {
	return "survey_votes"
}
