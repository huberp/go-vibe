package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"myapp/internal/models"
	"myapp/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	t.Run("should return all users", func(t *testing.T) {
		users := []models.User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "user"},
			{ID: 2, Name: "Bob", Email: "bob@example.com", Role: "admin"},
		}

		mockRepo.EXPECT().FindAll(gomock.Any()).Return(users, nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.GET("/users", handler.GetUsers)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, "Alice", response[0].Name)
	})

	t.Run("should handle database error", func(t *testing.T) {
		mockRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("database error"))

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.GET("/users", handler.GetUsers)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	t.Run("should create user successfully", func(t *testing.T) {
		userInput := CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
			Role:     "user",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, user *models.User) error {
				user.ID = 1
				return nil
			},
		)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should reject invalid input", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject missing required fields - name", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		userInput := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject missing required fields - email", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		userInput := map[string]string{
			"name":     "John Doe",
			"password": "password123",
		}

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject missing required fields - password", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		userInput := map[string]string{
			"name":  "John Doe",
			"email": "test@example.com",
		}

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject invalid email format", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		userInput := CreateUserRequest{
			Name:     "John Doe",
			Email:    "not-an-email",
			Password: "password123",
		}

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject short password", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		userInput := CreateUserRequest{
			Name:     "John Doe",
			Email:    "test@example.com",
			Password: "short",
		}

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject invalid role", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		userInput := CreateUserRequest{
			Name:     "John Doe",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "superadmin",
		}

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should use default role when not provided", func(t *testing.T) {
		userInput := CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, user *models.User) error {
				assert.Equal(t, "user", user.Role)
				user.ID = 1
				return nil
			},
		)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should create admin user when role is provided", func(t *testing.T) {
		userInput := CreateUserRequest{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "password123",
			Role:     "admin",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, user *models.User) error {
				assert.Equal(t, "admin", user.Role)
				user.ID = 1
				return nil
			},
		)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should handle database error on create", func(t *testing.T) {
		userInput := CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.POST("/users", handler.CreateUser)

		body, _ := json.Marshal(userInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	t.Run("should return user by ID", func(t *testing.T) {
		user := &models.User{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "user"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(user, nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Alice", response.Name)
	})

	t.Run("should return 404 when user not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(999)).Return(nil, repository.ErrUserNotFound)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(999))
			c.Set("user_role", "user")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should allow admin to access any user", func(t *testing.T) {
		user := &models.User{ID: 5, Name: "Bob", Email: "bob@example.com", Role: "user"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(5)).Return(user, nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/5", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Bob", response.Name)
	})

	t.Run("should allow user to access their own profile", func(t *testing.T) {
		user := &models.User{ID: 3, Name: "Charlie", Email: "charlie@example.com", Role: "user"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(3)).Return(user, nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(3))
			c.Set("user_role", "user")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/3", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Charlie", response.Name)
	})

	t.Run("should deny user from accessing another user's profile", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(3))
			c.Set("user_role", "user")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/10", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "insufficient permissions", response["error"])
	})

	t.Run("should return 400 for invalid ID format", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("should return 400 for negative ID", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/-1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should handle database error on FindByID", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(5)).Return(nil, errors.New("database connection error"))

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.GET("/users/:id", handler.GetUserByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/5", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "failed to fetch user", response["error"])
	})
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	t.Run("should update user successfully", func(t *testing.T) {
		updateInput := UpdateUserRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		existingUser := &models.User{ID: 1, Name: "Old Name", Email: "old@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(existingUser, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should allow admin to update any user", func(t *testing.T) {
		updateInput := UpdateUserRequest{
			Name:  "Admin Updated",
			Email: "admin.updated@example.com",
		}

		existingUser := &models.User{ID: 5, Name: "Original", Email: "original@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(5)).Return(existingUser, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/5", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should allow user to update their own profile", func(t *testing.T) {
		updateInput := UpdateUserRequest{
			Name:  "Self Updated",
			Email: "self.updated@example.com",
		}

		existingUser := &models.User{ID: 7, Name: "Original", Email: "original@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(7)).Return(existingUser, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(7))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/7", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should deny user from updating another user's profile", func(t *testing.T) {
		updateInput := UpdateUserRequest{
			Name:  "Unauthorized Update",
			Email: "unauthorized@example.com",
		}

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(7))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/10", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "insufficient permissions", response["error"])
	})

	t.Run("should return 400 for invalid ID format", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Name: "Test",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/invalid", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("should return 404 when user not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(999)).Return(nil, repository.ErrUserNotFound)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(999))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Name: "Test",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/999", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "user not found", response["error"])
	})

	t.Run("should handle database error on FindByID", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(5)).Return(nil, errors.New("database connection error"))

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "admin")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Name: "Test",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/5", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "failed to fetch user", response["error"])
	})

	t.Run("should reject invalid JSON", func(t *testing.T) {
		existingUser := &models.User{ID: 1, Name: "Test", Email: "test@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(existingUser, nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should reject invalid email format", func(t *testing.T) {
		existingUser := &models.User{ID: 1, Name: "Test", Email: "test@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(existingUser, nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Email: "not-an-email",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should update only name when email not provided", func(t *testing.T) {
		existingUser := &models.User{ID: 1, Name: "Old Name", Email: "old@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(existingUser, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, user *models.User) error {
				assert.Equal(t, "New Name", user.Name)
				assert.Equal(t, "old@example.com", user.Email)
				return nil
			},
		)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Name: "New Name",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should update only email when name not provided", func(t *testing.T) {
		existingUser := &models.User{ID: 1, Name: "Old Name", Email: "old@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(existingUser, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, user *models.User) error {
				assert.Equal(t, "Old Name", user.Name)
				assert.Equal(t, "new@example.com", user.Email)
				return nil
			},
		)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Email: "new@example.com",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle database error on Update", func(t *testing.T) {
		existingUser := &models.User{ID: 1, Name: "Old Name", Email: "old@example.com"}
		mockRepo.EXPECT().FindByID(gomock.Any(), uint(1)).Return(existingUser, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateUser)

		updateInput := UpdateUserRequest{
			Name: "New Name",
		}
		body, _ := json.Marshal(updateInput)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "failed to update user", response["error"])
	})
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	t.Run("should delete user successfully", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), uint(1)).Return(nil)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return 404 when user not found", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), uint(999)).Return(repository.ErrUserNotFound)

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid ID format", func(t *testing.T) {
		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "invalid user ID", response["error"])
	})

	t.Run("should handle database error", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), uint(5)).Return(errors.New("database connection error"))

		handler := NewUserHandler(mockRepo)
		router := gin.New()
		router.DELETE("/users/:id", handler.DeleteUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/5", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "failed to delete user", response["error"])
	})
}
