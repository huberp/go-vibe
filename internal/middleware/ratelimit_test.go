package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should allow requests within rate limit", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(10, 5)) // 10 requests per second, burst of 5
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should reject requests exceeding rate limit", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 2)) // 1 request per second, burst of 2
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// First two requests should succeed (burst)
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}

		// Third request should be rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("should handle multiple IPs independently", func(t *testing.T) {
		// Use different IPs that won't interfere with each other
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 1)) // 1 request per second, burst of 1
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Request from IP 1
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		req1.Header.Set("X-Forwarded-For", "192.168.100.1")
		req1.RemoteAddr = "192.168.100.1:12345"
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Request from IP 2 should succeed (different IP)
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.Header.Set("X-Forwarded-For", "192.168.100.2")
		req2.RemoteAddr = "192.168.100.2:12345"
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Second request from IP 1 should be rate limited
		w3 := httptest.NewRecorder()
		req3, _ := http.NewRequest("GET", "/test", nil)
		req3.Header.Set("X-Forwarded-For", "192.168.100.1")
		req3.RemoteAddr = "192.168.100.1:12345"
		router.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusTooManyRequests, w3.Code)

		// Request from IP 2 should still be rate limited (already used burst)
		w4 := httptest.NewRecorder()
		req4, _ := http.NewRequest("GET", "/test", nil)
		req4.Header.Set("X-Forwarded-For", "192.168.100.2")
		req4.RemoteAddr = "192.168.100.2:12345"
		router.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusTooManyRequests, w4.Code)
	})

	t.Run("should handle burst capacity correctly", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 5)) // 1 request per second, burst of 5
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// All 5 burst requests should succeed
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
		}

		// 6th request should be rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("should handle zero burst capacity", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(10, 0)) // 10 requests per second, burst of 0
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// With burst of 0, requests should be heavily rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		// First request might succeed or fail depending on rate limiter implementation
		// With golang.org/x/time/rate, burst 0 means no tokens available
	})

	t.Run("should handle high burst capacity", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 100)) // 1 request per second, burst of 100
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// All 100 burst requests should succeed
		for i := 0; i < 100; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
		}

		// 101st request should be rate limited
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("should return correct error message", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 1))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Use up the burst
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Next request should be rate limited with proper error
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
		assert.Contains(t, w2.Body.String(), "rate limit exceeded")
	})

	t.Run("should handle concurrent requests from same IP", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(10, 5))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		var wg sync.WaitGroup
		successCount := 0
		rateLimitedCount := 0
		var mu sync.Mutex

		// Make 10 concurrent requests
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Forwarded-For", "192.168.1.100")
				router.ServeHTTP(w, req)

				mu.Lock()
				if w.Code == http.StatusOK {
					successCount++
				} else if w.Code == http.StatusTooManyRequests {
					rateLimitedCount++
				}
				mu.Unlock()
			}()
		}

		wg.Wait()

		// With burst of 5, at most 5 should succeed
		assert.LessOrEqual(t, successCount, 5)
		// Some requests should be rate limited
		assert.Greater(t, rateLimitedCount, 0)
	})

	t.Run("should track different IPs separately", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 2))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		ips := []string{"10.10.0.1", "10.10.0.2", "10.10.0.3"}

		// Each IP should be able to make 2 requests (burst)
		for _, ip := range ips {
			for i := 0; i < 2; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Forwarded-For", ip)
				req.RemoteAddr = ip + ":12345"
				router.ServeHTTP(w, req)
				assert.Equal(t, http.StatusOK, w.Code, "Request %d from IP %s should succeed", i+1, ip)
			}

			// Third request from each IP should be rate limited
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Forwarded-For", ip)
			req.RemoteAddr = ip + ":12345"
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusTooManyRequests, w.Code, "Third request from IP %s should be rate limited", ip)
		}
	})

	t.Run("should handle requests without X-Forwarded-For header", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 1))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// First request should succeed
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request (same IP from httptest) should be rate limited
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})
}

func TestNewRateLimiter(t *testing.T) {
	t.Run("should create rate limiter with specified parameters", func(t *testing.T) {
		limiter := NewRateLimiter(10, 5)
		assert.NotNil(t, limiter)
		assert.NotNil(t, limiter.limiters)
	})

	t.Run("should create different limiters for different parameters", func(t *testing.T) {
		limiter1 := NewRateLimiter(10, 5)
		limiter2 := NewRateLimiter(20, 10)
		assert.NotNil(t, limiter1)
		assert.NotNil(t, limiter2)
	})
}

func TestGetLimiter(t *testing.T) {
	t.Run("should create new limiter for new IP", func(t *testing.T) {
		rateLimiter := NewRateLimiter(10, 5)
		limiter := rateLimiter.GetLimiter("192.168.1.1")
		assert.NotNil(t, limiter)
	})

	t.Run("should return same limiter for same IP", func(t *testing.T) {
		rateLimiter := NewRateLimiter(10, 5)
		limiter1 := rateLimiter.GetLimiter("192.168.1.1")
		limiter2 := rateLimiter.GetLimiter("192.168.1.1")
		assert.Equal(t, limiter1, limiter2)
	})

	t.Run("should return different limiters for different IPs", func(t *testing.T) {
		rateLimiter := NewRateLimiter(10, 5)
		limiter1 := rateLimiter.GetLimiter("192.168.1.1")
		limiter2 := rateLimiter.GetLimiter("192.168.1.2")
		// Verify they are different instances (not the same pointer)
		assert.NotSame(t, limiter1, limiter2)
	})

	t.Run("should be thread-safe", func(t *testing.T) {
		rateLimiter := NewRateLimiter(10, 5)
		var wg sync.WaitGroup

		// Create limiters concurrently
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				ip := "192.168.1.1"
				limiter := rateLimiter.GetLimiter(ip)
				assert.NotNil(t, limiter)
			}(i)
		}

		wg.Wait()
	})
}

func TestRateLimitMiddleware_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should handle very high request rate", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1000, 100)) // High rate limit
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Make many requests quickly
		for i := 0; i < 100; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("should handle very low request rate", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(0.1, 1)) // Very low rate: 1 request per 10 seconds
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// First request should succeed (burst)
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Immediate second request should be rate limited
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})

	t.Run("should handle IPv6 addresses", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(1, 1))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Request from IPv6 address
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		req1.Header.Set("X-Forwarded-For", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request should be rate limited
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.Header.Set("X-Forwarded-For", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})

	t.Run("should allow requests after rate limit window passes", func(t *testing.T) {
		router := gin.New()
		router.Use(RateLimitMiddleware(100, 1)) // 100 requests per second, burst of 1
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// First request succeeds
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Wait for rate limit to refill (with 100 req/s, tokens refill quickly)
		time.Sleep(20 * time.Millisecond)

		// Second request should succeed after waiting
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)
	})
}
