package models

import (
	"time"

	"gorm.io/gorm"
)

type Elderly struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"not null" binding:"required"`
	Gender           string         `json:"gender" gorm:"not null" binding:"required"` // 男, 女
	Age              int            `json:"age" gorm:"not null" binding:"required"`
	IDCard           *string        `json:"id_card" gorm:"uniqueIndex;type:varchar(18)"`
	Phone            string         `json:"phone"`
	Address          string         `json:"address"`
	HealthCondition  string         `json:"health_condition"` // 健康状况
	EconomicSource   string         `json:"economic_source"`  // 经济来源
	LivingCondition  string         `json:"living_condition"` // 居住情况
	CareLevel        string         `json:"care_level"`       // 护理等级
	MedicalInsurance string         `json:"medical_insurance"` // 医疗保险
	Remarks          string         `json:"remarks"`
	RecorderID       uint           `json:"recorder_id"`
	Recorder         *User          `json:"recorder,omitempty" gorm:"foreignKey:RecorderID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Elderly) TableName() string {
	return "elderly"
}
