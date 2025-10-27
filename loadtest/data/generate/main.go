package main

import (
	"flag"
	"fmt"
	"log"
	"myapp/internal/models"
	"myapp/pkg/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	var (
		dbURL     string
		numUsers  int
		adminUser bool
	)

	flag.StringVar(&dbURL, "db", "", "Database connection URL (required)")
	flag.IntVar(&numUsers, "users", 100, "Number of users to generate")
	flag.BoolVar(&adminUser, "admin", true, "Include admin user")
	flag.Parse()

	if dbURL == "" {
		log.Fatal("Database URL is required. Use -db flag.")
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Printf("Generating %d test users...\n", numUsers)

	// Generate admin user if requested
	if adminUser {
		hashedPassword, err := utils.HashPassword("password123")
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}
		admin := &models.User{
			Name:         "Load Test Admin",
			Email:        "loadtest-admin@example.com",
			PasswordHash: hashedPassword,
			Role:         "admin",
		}

		// Check if admin already exists
		var existingAdmin models.User
		result := db.Where("email = ?", admin.Email).First(&existingAdmin)
		if result.Error == nil {
			fmt.Println("Admin user already exists, skipping...")
		} else {
			if err := db.Create(admin).Error; err != nil {
				log.Printf("Warning: Failed to create admin user: %v", err)
			} else {
				fmt.Printf("Created admin user: %s\n", admin.Email)
			}
		}
	}

	// Generate regular users
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	batchSize := 100
	for i := 0; i < numUsers; i += batchSize {
		end := i + batchSize
		if end > numUsers {
			end = numUsers
		}

		var users []models.User
		for j := i; j < end; j++ {
			user := models.User{
				Name:         fmt.Sprintf("Load Test User %d", j+1),
				Email:        fmt.Sprintf("loadtest-user-%d@example.com", j+1),
				PasswordHash: hashedPassword,
				Role:         "user",
			}
			users = append(users, user)
		}

		// Batch insert
		if err := db.CreateInBatches(users, batchSize).Error; err != nil {
			log.Printf("Warning: Failed to create users batch %d-%d: %v", i, end, err)
		} else {
			fmt.Printf("Created users %d-%d\n", i+1, end)
		}
	}

	// Count total users
	var totalUsers int64
	db.Model(&models.User{}).Count(&totalUsers)
	fmt.Printf("\nTotal users in database: %d\n", totalUsers)

	// Count by role
	var adminCount, userCount int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	db.Model(&models.User{}).Where("role = ?", "user").Count(&userCount)
	fmt.Printf("Admin users: %d, Regular users: %d\n", adminCount, userCount)
}
