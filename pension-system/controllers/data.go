package controllers

import (
	"encoding/json"
	"pension-system/database"
	"pension-system/models"

	"gorm.io/gorm"
)

type DataController struct {
	db *gorm.DB
}

func NewDataController() *DataController {
	return &DataController{
		db: database.GetDB(),
	}
}

type ElderlyResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *models.Elderly `json:"data,omitempty"`
}

type ElderlyListResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []models.Elderly `json:"data,omitempty"`
	Total   int64          `json:"total"`
}

type SummaryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *SummaryData `json:"data,omitempty"`
}

type SummaryData struct {
	TotalElderly      int64               `json:"total_elderly"`
	GenderDistribution map[string]int64   `json:"gender_distribution"`
	AgeDistribution   map[string]int64   `json:"age_distribution"`
	HealthConditions  map[string]int64   `json:"health_conditions"`
	CareLevels        map[string]int64   `json:"care_levels"`
}

// GetElderlyList returns all elderly records
func (c *DataController) GetElderlyList(page, pageSize int, keyword string) (string, error) {
	var elderly []models.Elderly
	var total int64

	query := c.db.Model(&models.Elderly{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR id_card LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := query.Preload("Recorder").Offset(offset).Limit(pageSize).Find(&elderly)

	if result.Error != nil {
		resp := ElderlyListResponse{
			Success: false,
			Message: "获取数据失败",
			Total:   0,
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := ElderlyListResponse{
		Success: true,
		Message: "获取成功",
		Data:    elderly,
		Total:   total,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetElderly returns a single elderly record by ID
func (c *DataController) GetElderly(id uint) (string, error) {
	var elderly models.Elderly
	result := c.db.Preload("Recorder").First(&elderly, id)

	if result.Error != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "获取数据失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := ElderlyResponse{
		Success: true,
		Message: "获取成功",
		Data:    &elderly,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CreateElderly creates a new elderly record
func (c *DataController) CreateElderly(dataStr string) (string, error) {
	var elderly models.Elderly
	err := json.Unmarshal([]byte(dataStr), &elderly)
	if err != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "数据格式错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), err
	}

	result := c.db.Create(&elderly)
	if result.Error != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "创建失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Reload with associations
	c.db.Preload("Recorder").First(&elderly, elderly.ID)

	resp := ElderlyResponse{
		Success: true,
		Message: "创建成功",
		Data:    &elderly,
	}

	resultData, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(resultData), nil
}

// UpdateElderly updates an existing elderly record
func (c *DataController) UpdateElderly(id uint, dataStr string) (string, error) {
	var elderly models.Elderly
	result := c.db.First(&elderly, id)
	if result.Error != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "记录不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	var updateData models.Elderly
	err := json.Unmarshal([]byte(dataStr), &updateData)
	if err != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "数据格式错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), err
	}

	result = c.db.Model(&elderly).Updates(updateData)
	if result.Error != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "更新失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Reload with associations
	c.db.Preload("Recorder").First(&elderly, elderly.ID)

	resp := ElderlyResponse{
		Success: true,
		Message: "更新成功",
		Data:    &elderly,
	}

	resultData, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(resultData), nil
}

// DeleteElderly deletes an elderly record
func (c *DataController) DeleteElderly(id uint) (string, error) {
	result := c.db.Delete(&models.Elderly{}, id)
	if result.Error != nil {
		resp := ElderlyResponse{
			Success: false,
			Message: "删除失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := ElderlyResponse{
		Success: true,
		Message: "删除成功",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetSummary returns summary statistics
func (c *DataController) GetSummary() (string, error) {
	var totalElderly int64
	c.db.Model(&models.Elderly{}).Count(&totalElderly)

	// Gender distribution
	var genderData []struct {
		Gender string
		Count  int64
	}
	c.db.Model(&models.Elderly{}).Select("gender, count(*) as count").Group("gender").Scan(&genderData)
	genderDistribution := make(map[string]int64)
	for _, item := range genderData {
		genderDistribution[item.Gender] = item.Count
	}

	// Age distribution
	var ageData []struct {
		AgeGroup string
		Count    int64
	}
	c.db.Raw(`
		SELECT
			CASE
				WHEN age < 60 THEN '60岁以下'
				WHEN age BETWEEN 60 AND 69 THEN '60-69岁'
				WHEN age BETWEEN 70 AND 79 THEN '70-79岁'
				WHEN age BETWEEN 80 AND 89 THEN '80-89岁'
				ELSE '90岁以上'
			END as age_group,
			COUNT(*) as count
		FROM elderly
		GROUP BY age_group
	`).Scan(&ageData)

	ageDistribution := make(map[string]int64)
	for _, item := range ageData {
		ageDistribution[item.AgeGroup] = item.Count
	}

	// Health conditions
	var healthData []struct {
		Condition string
		Count     int64
	}
	c.db.Model(&models.Elderly{}).Select("health_condition, count(*) as count").
		Group("health_condition").Scan(&healthData)
	healthConditions := make(map[string]int64)
	for _, item := range healthData {
		if item.Condition != "" {
			healthConditions[item.Condition] = item.Count
		}
	}

	// Care levels
	var careData []struct {
		Level string
		Count int64
	}
	c.db.Model(&models.Elderly{}).Select("care_level, count(*) as count").
		Group("care_level").Scan(&careData)
	careLevels := make(map[string]int64)
	for _, item := range careData {
		if item.Level != "" {
			careLevels[item.Level] = item.Count
		}
	}

	summary := &SummaryData{
		TotalElderly:      totalElderly,
		GenderDistribution: genderDistribution,
		AgeDistribution:   ageDistribution,
		HealthConditions:  healthConditions,
		CareLevels:        careLevels,
	}

	resp := SummaryResponse{
		Success: true,
		Message: "获取成功",
		Data:    summary,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
