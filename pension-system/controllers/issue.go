package controllers

import (
	"encoding/json"
	"pension-system/database"
	"pension-system/models"

	"gorm.io/gorm"
)

type IssueController struct {
	db *gorm.DB
}

func NewIssueController() *IssueController {
	return &IssueController{
		db: database.GetDB(),
	}
}

type IssueResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *models.Issue `json:"data,omitempty"`
}

type IssueListResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []models.Issue `json:"data,omitempty"`
	Total   int64        `json:"total"`
}

// GetIssueList returns all issues
func (c *IssueController) GetIssueList(page, pageSize int, status string) (string, error) {
	var issues []models.Issue
	var total int64

	query := c.db.Model(&models.Issue{})

	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := c.db.Preload("Submitter").Preload("Handler").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&issues)

	if result.Error != nil {
		resp := IssueListResponse{
			Success: false,
			Message: "获取问题失败",
			Total:   0,
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := IssueListResponse{
		Success: true,
		Message: "获取成功",
		Data:    issues,
		Total:   total,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetIssue returns a single issue by ID
func (c *IssueController) GetIssue(id uint) (string, error) {
	var issue models.Issue
	result := c.db.Preload("Submitter").Preload("Handler").First(&issue, id)

	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "获取问题失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := IssueResponse{
		Success: true,
		Message: "获取成功",
		Data:    &issue,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CreateIssue creates a new issue
func (c *IssueController) CreateIssue(title, description, category, priority string, submitterID uint) (string, error) {
	issue := models.Issue{
		Title:       title,
		Description: description,
		Category:    category,
		Priority:    priority,
		Status:      "pending",
		SubmitterID: submitterID,
	}

	result := c.db.Create(&issue)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "提交失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	c.db.Preload("Submitter").Preload("Handler").First(&issue, issue.ID)

	resp := IssueResponse{
		Success: true,
		Message: "提交成功",
		Data:    &issue,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// UpdateIssueStatus updates issue status
func (c *IssueController) UpdateIssueStatus(id uint, status string) (string, error) {
	var issue models.Issue
	result := c.db.First(&issue, id)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "问题不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	result = c.db.Model(&issue).Update("status", status)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "更新失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	c.db.Preload("Submitter").Preload("Handler").First(&issue, issue.ID)

	resp := IssueResponse{
		Success: true,
		Message: "更新成功",
		Data:    &issue,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// RespondToIssue adds a response to an issue
func (c *IssueController) RespondToIssue(id uint, response string, handlerID uint) (string, error) {
	var issue models.Issue
	result := c.db.First(&issue, id)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "问题不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	updates := map[string]interface{}{
		"response":  response,
		"handler_id": handlerID,
		"status":    "processing",
	}

	result = c.db.Model(&issue).Updates(updates)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "回复失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	c.db.Preload("Submitter").Preload("Handler").First(&issue, issue.ID)

	resp := IssueResponse{
		Success: true,
		Message: "回复成功",
		Data:    &issue,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CloseIssue closes an issue
func (c *IssueController) CloseIssue(id uint) (string, error) {
	var issue models.Issue
	result := c.db.First(&issue, id)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "问题不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	result = c.db.Model(&issue).Update("status", "closed")
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "关闭失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	c.db.Preload("Submitter").Preload("Handler").First(&issue, issue.ID)

	resp := IssueResponse{
		Success: true,
		Message: "已关闭",
		Data:    &issue,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DeleteIssue deletes an issue
func (c *IssueController) DeleteIssue(id uint) (string, error) {
	result := c.db.Delete(&models.Issue{}, id)
	if result.Error != nil {
		resp := IssueResponse{
			Success: false,
			Message: "删除失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := IssueResponse{
		Success: true,
		Message: "删除成功",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
