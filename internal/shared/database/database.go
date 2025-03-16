package database

import (
	"fmt"
	"log"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/shared/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database represents the database connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(config *config.Config, l *logger.ZapLogger) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.DB.HOST, config.DB.USER, config.DB.PASSWORD, config.DB.NAME, config.DB.PORT, config.DB.SSL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	l.Info("Connected to database")

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{
		DB: db,
	}, nil
}

// GetDB returns the GORM database instance
func (d *Database) GetDB() *gorm.DB {
	return d.DB
}
