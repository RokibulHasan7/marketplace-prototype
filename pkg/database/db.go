package database

import (
	"fmt"
	"github.com/RokibulHasan7/marketplace-prototype/pkg/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL") // Example: "postgres://user:password@localhost:5433/marketplace_db"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to the database:", err)
	}

	// Auto Migrate Tables
	err = db.AutoMigrate(&models.User{}, &models.Application{}, &models.Deployment{}, &models.BillingRecord{})
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	DB = db
	fmt.Println("✅ Database connected & migrated successfully!")
}
