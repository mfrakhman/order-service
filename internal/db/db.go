package db

import (
	"fmt"
	"log"
	"ordersvc/config"
	"ordersvc/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("❌ Failed to migrate: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL and migrated")
	return db
}
