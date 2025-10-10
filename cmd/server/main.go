package main

import (
	"fmt"
	"log"
	"myapp/internal/models"
	"myapp/internal/routes"
	"myapp/pkg/config"
	"myapp/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger
	if err := logger.Initialize(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(router, db, logger.Log, cfg.JWTSecret)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Log.Info("Starting server", zap.String("addr", addr))

	if err := router.Run(addr); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
