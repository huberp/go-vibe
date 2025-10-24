package routes

import (
	"myapp/internal/handlers"
	"myapp/internal/middleware"
	"myapp/internal/repository"
	"myapp/pkg/config"
	"myapp/pkg/health"
	"myapp/pkg/info"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, logger *zap.Logger, jwtSecret string) {
	// Load configuration
	cfg := config.Load()

	// Configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database connection", zap.Error(err))
	}

	// Apply settings from Viper configuration
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Create repository
	userRepo := repository.NewPostgresUserRepository(db)

	// Create handlers
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(db, jwtSecret, logger)
	
	// Setup health check providers
	healthRegistry := health.NewRegistry()
	healthRegistry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady))
	healthHandler := handlers.NewHealthHandler(healthRegistry)

	// Setup info providers
	infoRegistry := info.NewRegistry()
	infoRegistry.Register(info.NewBuildInfoProvider("dev", "unknown", "", runtime.Version()))
	infoRegistry.Register(info.NewUserStatsProvider(db))
	infoHandler := handlers.NewInfoHandler(infoRegistry)

	// Register user count metric collector
	middleware.RegisterUserCountCollector(db)

	// CORS middleware - allow all origins in development, configure for production
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure this for production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "traceparent", "tracestate"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// OpenTelemetry tracing middleware (conditionally enabled)
	router.Use(middleware.OtelMiddleware("myapp", cfg.Observability.Otel))

	// Logging middleware (with W3C trace context support)
	router.Use(middleware.LoggingMiddleware(logger))

	// Rate limiting middleware - configurable via YAML/environment variables
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst))

	// Prometheus metrics middleware
	router.Use(middleware.PrometheusMiddleware())

	// Recovery middleware
	router.Use(gin.Recovery())

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoints
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/startup", healthHandler.StartupProbe)
	router.GET("/health/liveness", healthHandler.LivenessProbe)
	router.GET("/health/readiness", healthHandler.ReadinessProbe)

	// Info endpoint
	router.GET("/info", infoHandler.GetInfo)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Public routes
		v1.POST("/login", authHandler.Login)
		v1.POST("/users", userHandler.CreateUser) // Public signup

		// Protected routes
		protected := v1.Group("/")
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

	// Legacy routes (backward compatibility) - redirect to v1
	router.POST("/login", authHandler.Login)
	router.POST("/users", userHandler.CreateUser)

	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	{
		admin := protected.Group("/")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/users", userHandler.GetUsers)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
		}

		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
	}
}
