package main

import (
	"flag"
	"fmt"
	"log"
	"myapp/internal/models"
	"myapp/internal/routes"
	"myapp/pkg/config"
	"myapp/pkg/logger"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Parse command line flags
	var stage string
	flag.StringVar(&stage, "stage", "", "Configuration stage (development, staging, production)")
	flag.Parse()

	// Set stage environment variable if flag is provided
	if stage != "" {
		os.Setenv("APP_STAGE", stage)
	}

	// Initialize logger
	if err := logger.Initialize(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Log the active stage
	activeStage := os.Getenv("APP_STAGE")
	if activeStage == "" {
		activeStage = "development"
	}
	logger.Log.Info("Starting server", 
		zap.String("stage", activeStage),
		zap.String("port", cfg.Server.Port))

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
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
	routes.SetupRoutes(router, db, logger.Log, cfg.JWT.Secret)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Log.Info("Server listening", zap.String("addr", addr))

	if err := router.Run(addr); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
