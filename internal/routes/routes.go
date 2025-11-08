package routes

import (
	"myapp/internal/handlers"
	"myapp/internal/middleware"
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

	// Add your domain-specific repositories and handlers here
	// Example:
	// myRepo := repository.NewMyRepository(db)
	// myHandler := handlers.NewMyHandler(myRepo)

	// Setup health check providers
	healthRegistry := health.NewRegistry()
	healthRegistry.Register(health.NewDatabaseHealthCheckProvider(db, health.ScopeStartup, health.ScopeReady))
	healthHandler := handlers.NewHealthHandler(healthRegistry)

	// Setup info providers
	infoRegistry := info.NewRegistry()
	infoRegistry.Register(info.NewBuildInfoProvider("dev", "unknown", "", runtime.Version()))
	// Add your custom info providers here
	// Example: infoRegistry.Register(info.NewMyStatsProvider(db))
	infoHandler := handlers.NewInfoHandler(infoRegistry)

	// Register custom metric collectors here
	// Example: middleware.RegisterMyCountCollector(db)

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

	// Add your API routes here
	// Example:
	// v1 := router.Group("/v1")
	// {
	//     // Public routes
	//     v1.POST("/login", authHandler.Login)
	//     v1.POST("/items", myHandler.CreateItem)
	//
	//     // Protected routes
	//     protected := v1.Group("/")
	//     protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
	//     {
	//         protected.GET("/items", myHandler.GetItems)
	//         protected.GET("/items/:id", myHandler.GetItemByID)
	//         protected.PUT("/items/:id", myHandler.UpdateItem)
	//         protected.DELETE("/items/:id", myHandler.DeleteItem)
	//     }
	// }
	//
	// See examples/user-management for a complete working example
}
