package routes

import (
	"myapp/internal/handlers"
	"myapp/internal/middleware"
	"myapp/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, jwtSecret string) {
	// Create repository
	userRepo := repository.NewPostgresUserRepository(db)

	// Create handlers
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(db, jwtSecret)

	// Middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.PrometheusMiddleware())
	router.Use(gin.Recovery())

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Public routes
	router.POST("/login", authHandler.Login)
	router.POST("/users", userHandler.CreateUser) // Public signup

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		// Admin-only routes
		admin := protected.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", userHandler.GetUsers)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
		}

		// Owner or admin routes
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
	}
}
