package handlers

import (
	"myapp/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	db     *gorm.DB
	secret string
	logger *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB, secret string, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		secret: secret,
		logger: logger,
	}
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

// Login authenticates a user and returns a JWT token
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Router /v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid login request",
			zap.String("error", err.Error()),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get request context for logging
	requestID, _ := c.Get("request_id")
	clientIP := c.ClientIP()

	// Find user by email
	var user struct {
		ID           uint
		Name         string
		Email        string
		PasswordHash string
		Role         string
	}

	if err := h.db.WithContext(c.Request.Context()).Table("users").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.logger.Warn("login attempt with unknown email",
				zap.String("email", req.Email),
				zap.String("client_ip", clientIP),
				zap.Any("request_id", requestID),
			)
		} else {
			h.logger.Error("database error during login",
				zap.Error(err),
				zap.String("email", req.Email),
				zap.String("client_ip", clientIP),
				zap.Any("request_id", requestID),
			)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		h.logger.Warn("login attempt with invalid password",
			zap.String("email", req.Email),
			zap.String("client_ip", clientIP),
			zap.Any("request_id", requestID),
			zap.Uint("user_id", user.ID),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Role, h.secret)
	if err != nil {
		h.logger.Error("failed to generate JWT token",
			zap.Error(err),
			zap.Uint("user_id", user.ID),
			zap.String("email", req.Email),
			zap.String("client_ip", clientIP),
			zap.Any("request_id", requestID),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Log successful authentication
	h.logger.Info("successful login",
		zap.Uint("user_id", user.ID),
		zap.String("email", req.Email),
		zap.String("role", user.Role),
		zap.String("client_ip", clientIP),
		zap.Any("request_id", requestID),
	)

	response := LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Name = user.Name
	response.User.Email = user.Email
	response.User.Role = user.Role

	c.JSON(http.StatusOK, response)
}
