package main

import (
	"embed"
	"pension-system/controllers"
	"pension-system/database"
	"pension-system/models"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Initialize database before creating controllers
	db, err := database.InitDB()
	if err != nil {
		println("Failed to initialize database:", err.Error())
		return
	}

	// Auto migrate tables
	err = database.AutoMigrate(db, &models.User{}, &models.Elderly{}, &models.Survey{}, &models.Issue{})
	if err != nil {
		println("Failed to migrate database:", err.Error())
		return
	}

	// Create default admin user if not exists
	var adminCount int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		admin := &models.User{
			Username: "admin",
			Password: "$2a$10$iuy7mUFrqmFHswp.SvQuze7UtY/sl0qmZsTDoM3IcNclMG4E8e3dC", // admin123
			Role:     "admin",
			Name:     "系统管理员",
		}
		db.Create(admin)
	}

	// Create controllers with initialized database
	authController := controllers.NewAuthController()
	dataController := controllers.NewDataController()
	surveyController := controllers.NewSurveyController()
	issueController := controllers.NewIssueController()
	accountController := controllers.NewAccountController()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "智慧养老系统",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		OnDomReady:       app.domReady,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
			authController,
			dataController,
			surveyController,
			issueController,
			accountController,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
