package database

import (
	"fmt"
	"log"
	"os"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndAutoMigrate() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	// Run auto migrations
	err = db.AutoMigrate(
		models.Coffee{},
		models.User{},
		models.Order{},
		models.OrderItem{},
		models.Transaction{},
	)
	if err != nil {
		log.Fatalf("failed to run auto migrations: %v", err)
		return nil, err
	}

	return db, nil
}
