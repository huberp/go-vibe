package main

import (
	"flag"
	"fmt"
	"log"
	"myapp/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	var dbURL string
	flag.StringVar(&dbURL, "db", "", "Database connection URL (required)")
	flag.Parse()

	if dbURL == "" {
		log.Fatal("Database URL is required. Use -db flag.")
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Cleaning up load test data...")

	// Delete load test users
	result := db.Where("email LIKE ?", "loadtest-%@example.com").Delete(&models.User{})
	if result.Error != nil {
		log.Fatalf("Failed to clean up load test data: %v", result.Error)
	}

	fmt.Printf("Deleted %d load test users\n", result.RowsAffected)

	// Count remaining users
	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)
	fmt.Printf("Remaining users in database: %d\n", totalUsers)
}
