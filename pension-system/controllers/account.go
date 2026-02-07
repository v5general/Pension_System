package controllers

import (
	"encoding/json"
	"pension-system/database"
	"pension-system/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountController struct {
	db *gorm.DB
}

func NewAccountController() *AccountController {
	return &AccountController{
		db: database.GetDB(),
	}
}

type AccountResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    *models.User `json:"user,omitempty"`
}

// UpdateName updates user's display name
func (c *AccountController) UpdateName(newName, password string, userId uint) (string, error) {
	var user models.User
	result := c.db.First(&user, userId)
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "用户不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		resp := AccountResponse{
			Success: false,
			Message: "当前密码错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Check if new name is already taken
	var existingUser models.User
	c.db.Where("username = ? AND id != ?", newName, userId).First(&existingUser)
	if existingUser.ID != 0 {
		resp := AccountResponse{
			Success: false,
			Message: "该用户名已被使用",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Update name
	result = c.db.Model(&user).Update("name", newName)
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "修改失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Refresh user data
	c.db.First(&user, userId)

	resp := AccountResponse{
		Success: true,
		Message: "用户名修改成功",
		User:    &user,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// UpdatePassword updates user's password
func (c *AccountController) UpdatePassword(currentPassword, newPassword string, userId uint) (string, error) {
	var user models.User
	result := c.db.First(&user, userId)
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "用户不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		resp := AccountResponse{
			Success: false,
			Message: "当前密码错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		resp := AccountResponse{
			Success: false,
			Message: "密码加密失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), err
	}

	// Update password
	result = c.db.Model(&user).Update("password", string(hashedPassword))
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "密码修改失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := AccountResponse{
		Success: true,
		Message: "密码修改成功",
	}

	resultData, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(resultData), nil
}

// DeleteAccount deletes user's account
func (c *AccountController) DeleteAccount(password string, userId uint) (string, error) {
	var user models.User
	result := c.db.First(&user, userId)
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "用户不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		resp := AccountResponse{
			Success: false,
			Message: "当前密码错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Check if user is the last admin
	if user.Role == "admin" {
		var adminCount int64
		c.db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount <= 1 {
			resp := AccountResponse{
				Success: false,
				Message: "系统至少需要保留一个管理员账号，请先提升其他用户为管理员",
			}
			data, _ := json.Marshal(resp)
			return string(data), nil
		}
	}

	// Start transaction
	tx := c.db.Begin()

	// Delete user's votes
	tx.Where("user_id = ?", userId).Delete(&models.SurveyVote{})

	// Delete user's issues (soft delete)
	tx.Where("submitter_id = ?", userId).Delete(&models.Issue{})

	// Delete user account
	tx.Delete(&user)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		resp := AccountResponse{
			Success: false,
			Message: "注销失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), err
	}

	resp := AccountResponse{
		Success: true,
		Message: "账号已注销",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetUserList returns all users with their status (for admin user management)
func (c *AccountController) GetUserList(page, pageSize int) (string, error) {
	var users []models.User
	var total int64

	c.db.Model(&models.User{}).Count(&total)

	offset := (page - 1) * pageSize
	result := c.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	if result.Error != nil {
		resp := struct {
			Success bool          `json:"success"`
			Message string        `json:"message"`
			Data    []models.User `json:"data,omitempty"`
			Total   int64         `json:"total"`
		}{
			Success: false,
			Message: "获取用户列表失败",
			Total:   0,
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Calculate online status (simplified - last activity within 30 minutes)
	now := time.Now()
	for i := range users {
		users[i].LastSeen = now // In real implementation, track last login/activity time
	}

	resp := struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Data    []models.User `json:"data,omitempty"`
		Total   int64         `json:"total"`
	}{
		Success: true,
		Message: "获取成功",
		Data:    users,
		Total:   total,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DeleteUserAccount deletes a user account (admin function)
func (c *AccountController) DeleteUserAccount(targetUserId uint, adminId uint) (string, error) {
	// Get target user
	var targetUser models.User
	result := c.db.First(&targetUser, targetUserId)
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "用户不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Prevent deleting self
	if targetUserId == adminId {
		resp := AccountResponse{
			Success: false,
			Message: "不能注销自己的账号",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Check if target is the last admin
	if targetUser.Role == "admin" {
		var adminCount int64
		c.db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount <= 1 {
			resp := AccountResponse{
				Success: false,
				Message: "系统至少需要保留一个管理员账号",
			}
			data, _ := json.Marshal(resp)
			return string(data), nil
		}
	}

	// Start transaction
	tx := c.db.Begin()

	// Delete user's votes
	tx.Where("user_id = ?", targetUserId).Delete(&models.SurveyVote{})

	// Delete user's issues
	tx.Where("submitter_id = ?", targetUserId).Delete(&models.Issue{})

	// Delete user account
	tx.Delete(&targetUser)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		resp := AccountResponse{
			Success: false,
			Message: "注销失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), err
	}

	resp := AccountResponse{
		Success: true,
		Message: "用户账号已注销",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ResetUserPassword resets a user's password to default "123456" (admin function)
func (c *AccountController) ResetUserPassword(userId uint) (string, error) {
	// Get target user
	var user models.User
	result := c.db.First(&user, userId)
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "用户不存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	// Default password is 123456
	defaultPassword := "123456"

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		resp := AccountResponse{
			Success: false,
			Message: "密码加密失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), err
	}

	// Update password
	result = c.db.Model(&user).Update("password", string(hashedPassword))
	if result.Error != nil {
		resp := AccountResponse{
			Success: false,
			Message: "密码重置失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := AccountResponse{
		Success: true,
		Message: "密码已重置为：123456",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
