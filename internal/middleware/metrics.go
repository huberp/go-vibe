package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	userCountOnce sync.Once
)

// RegisterUserCountCollector registers a custom collector for user count
func RegisterUserCountCollector(db *gorm.DB) {
	userCountOnce.Do(func() {
		// Register a gauge that gets updated
		promauto.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "users_total",
				Help: "Total number of users in the database",
			},
			func() float64 {
				var count int64
				// Silently ignore errors to avoid panics in metric collection
				db.Table("users").Count(&count)
				return float64(count)
			},
		)
	})
}

// PrometheusMiddleware records metrics for HTTP requests
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.Request.URL.Path).Observe(duration)
	}
}
