package controllers

import (
	"encoding/json"
	"pension-system/database"
	"pension-system/models"
	"pension-system/utils"
	"strings"

	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController() *AuthController {
	return &AuthController{
		db: getDB(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"` // 这是账户名，需要唯一
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *UserResponse `json:"user,omitempty"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

// Login handles user login (use name/account for login)
func (c *AuthController) Login(name, password string) (string, error) {
	var user models.User
	result := c.db.Where("name = ?", name).First(&user)
	if result.Error != nil {
		resp := LoginResponse{
			Success: false,
			Message: "账户名或密码错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil // Return nil error so frontend can parse JSON
	}

	if !utils.CheckPassword(password, user.Password) {
		resp := LoginResponse{
			Success: false,
			Message: "账户名或密码错误",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	resp := LoginResponse{
		Success: true,
		Message: "登录成功",
		User: &UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Role:     user.Role,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Register handles user registration
func (c *AuthController) Register(username, password, name string) (string, error) {
	// Check if account name (name field) already exists (excluding soft-deleted records)
	var existingUser models.User
	result := c.db.Where("name = ?", name).First(&existingUser)
	// Not found is OK, other errors are not
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		resp := LoginResponse{
			Success: false,
			Message: "注册失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}
	// If found, it's a duplicate
	if result.Error == nil {
		resp := LoginResponse{
			Success: false,
			Message: "账户名已存在",
		}
		data, _ := json.Marshal(resp)
		return string(data), nil
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	// Create new user
	user := models.User{
		Username: username,
		Password: hashedPassword,
		Name:     name,
		Role:     "user",
	}
	result = c.db.Create(&user)
	if result.Error != nil {
		// Check if it's a unique constraint violation
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.name") ||
		   strings.Contains(result.Error.Error(), "unique constraint") {
			resp := LoginResponse{
				Success: false,
				Message: "账户名已存在",
			}
			data, _ := json.Marshal(resp)
			return string(data), nil
		}
		resp := LoginResponse{
			Success: false,
			Message: "注册失败",
		}
		data, _ := json.Marshal(resp)
		return string(data), result.Error
	}

	resp := LoginResponse{
		Success: true,
		Message: "注册成功",
		User: &UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Role:     user.Role,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GetCurrentUser returns the current logged in user
func (c *AuthController) GetCurrentUser(userID uint) (string, error) {
	var user models.User
	result := c.db.First(&user, userID)
	if result.Error != nil {
		return "", result.Error
	}

	resp := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getDB() *gorm.DB {
	return database.GetDB()
}
