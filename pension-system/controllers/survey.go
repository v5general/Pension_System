package controllers

import (
	"encoding/json"
	"pension-system/database"
	"pension-system/models"
	"time"

	"gorm.io/gorm"
)

type SurveyController struct {
	db *gorm.DB
}

func NewSurveyController() *SurveyController {
	return &SurveyController{
		db: database.GetDB(),
	}
}

type SurveyResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *models.Survey `json:"data,omitempty"`
}

type SurveyListResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []models.Survey `json:"data,omitempty"`
	Total   int64         `json:"total"`
}

// GetSurveyList returns all surveys
func (c *SurveyController) GetSurveyList(page, pageSize int) (string, error) {
	var surveys []models.Survey
	var total int64

	c.db.Model(&models.Survey{}).Count(&total)

	offset := (page - 1) * pageSize
	result := c.db.Preload("Creator").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&surveys)

	if result.Error != nil {
		resp := SurveyListResponse{
			Success: false,
			Message: "获取调查失败",
			Total:   0,
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := SurveyListResponse{
		Success: true,
		Message: "获取成功",
		Data:    surveys,
		Total:   total,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetActiveSurveys returns all active surveys
func (c *SurveyController) GetActiveSurveys() (string, error) {
	var surveys []models.Survey
	now := time.Now()

	result := c.db.Where("status = ? AND start_date <= ? AND end_date >= ?", "active", now, now).
		Preload("Creator").Find(&surveys)

	if result.Error != nil {
		resp := SurveyListResponse{
			Success: false,
			Message: "获取调查失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := SurveyListResponse{
		Success: true,
		Message: "获取成功",
		Data:    surveys,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetSurvey returns a single survey by ID
func (c *SurveyController) GetSurvey(id uint) (string, error) {
	var survey models.Survey
	result := c.db.Preload("Creator").First(&survey, id)

	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "获取调查失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := SurveyResponse{
		Success: true,
		Message: "获取成功",
		Data:    &survey,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CreateSurvey creates a new survey
func (c *SurveyController) CreateSurvey(title, description, options string, startDate, endDate string, createdBy uint) (string, error) {
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	survey := models.Survey{
		Title:       title,
		Description: description,
		StartDate:   start,
		EndDate:     end,
		Status:      "draft",
		Options:     options,
		CreatedBy:   createdBy,
	}

	result := c.db.Create(&survey)
	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "创建失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	c.db.Preload("Creator").First(&survey, survey.ID)

	resp := SurveyResponse{
		Success: true,
		Message: "创建成功",
		Data:    &survey,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// UpdateSurveyStatus updates survey status
func (c *SurveyController) UpdateSurveyStatus(id uint, status string) (string, error) {
	var survey models.Survey
	result := c.db.First(&survey, id)
	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "调查不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	result = c.db.Model(&survey).Update("status", status)
	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "更新失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	c.db.Preload("Creator").First(&survey, survey.ID)

	resp := SurveyResponse{
		Success: true,
		Message: "更新成功",
		Data:    &survey,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Vote submits a vote for a survey
func (c *SurveyController) Vote(surveyID uint, option string, userID uint) (string, error) {
	// Check if already voted
	var existingVote models.SurveyVote
	result := c.db.Where("survey_id = ? AND user_id = ?", surveyID, userID).First(&existingVote)
	if result.Error == nil {
		resp := SurveyResponse{
			Success: false,
			Message: "您已经参与过此调查",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Create vote
	vote := models.SurveyVote{
		SurveyID: surveyID,
		Option:   option,
		UserID:   userID,
	}

	result = c.db.Create(&vote)
	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "投票失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Update total votes
	c.db.Model(&models.Survey{}).Where("id = ?", surveyID).
		UpdateColumn("total_votes", gorm.Expr("total_votes + ?", 1))

	resp := SurveyResponse{
		Success: true,
		Message: "投票成功",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetSurveyResults returns survey results
func (c *SurveyController) GetSurveyResults(surveyID uint) (string, error) {
	var votes []models.SurveyVote
	result := c.db.Where("survey_id = ?", surveyID).Find(&votes)

	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "获取结果失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Count votes by option
	optionCounts := make(map[string]int)
	for _, vote := range votes {
		optionCounts[vote.Option]++
	}

	type ResultData struct {
		OptionCounts map[string]int `json:"option_counts"`
		TotalVotes   int            `json:"total_votes"`
	}

	resultData := ResultData{
		OptionCounts: optionCounts,
		TotalVotes:   len(votes),
	}

	data, err := json.Marshal(map[string]interface{}{
		"success": true,
		"message": "获取成功",
		"data":    resultData,
	})
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DeleteSurvey deletes a survey
func (c *SurveyController) DeleteSurvey(id uint) (string, error) {
	// Delete associated votes first
	c.db.Where("survey_id = ?", id).Delete(&models.SurveyVote{})

	// Delete survey
	result := c.db.Delete(&models.Survey{}, id)
	if result.Error != nil {
		resp := SurveyResponse{
			Success: false,
			Message: "删除失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := SurveyResponse{
		Success: true,
		Message: "删除成功",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
