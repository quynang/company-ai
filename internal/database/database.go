package database

import (
	"company-ai-training/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable vector extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return nil, err
	}

	// Auto migrate tables
	err = db.AutoMigrate(
		&models.Document{},
		&models.DocumentChunk{},
		&models.User{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.HRTicket{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
