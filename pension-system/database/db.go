package database

import (
	"database/sql"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	var err error
	sqlDB, err := sql.Open("sqlite", "pension_system.db")
	if err != nil {
		return nil, err
	}

	db, err = gorm.Open(sqlite.Dialector{
		Conn: sqlDB,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// AutoMigrate runs auto migration for given models
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

// Close closes the database connection
func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}
