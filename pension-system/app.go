package main

import (
	"context"
	"log"
	"pension-system/database"

	"gorm.io/gorm"
)

type App struct {
	ctx context.Context
	db  *gorm.DB
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Database is already initialized in main(), just get the instance
	a.db = database.GetDB()
	if a.db == nil {
		log.Fatal("Database not initialized")
	}

	log.Println("Application started successfully")
}

func (a *App) shutdown(ctx context.Context) {
	database.Close()
	log.Println("Application shutdown")
}

func (a *App) domReady(ctx context.Context) {
	// This function is called when the DOM is ready
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// This function is called before the application is closed
	// Return true to prevent the close, false to allow it
	return false
}

// GetDB returns the database instance
func (a *App) GetDB() *gorm.DB {
	return a.db
}
