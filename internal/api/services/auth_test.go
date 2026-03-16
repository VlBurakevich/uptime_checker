package services

import (
	"testing"
	"time"
	"uptime-checker/internal/api/dto"
	"uptime-checker/internal/api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestAuthDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Role{}, &models.User{}, &models.Credential{})
	require.NoError(t, err)

	db.Create(&models.Role{Name: models.RoleUser})

	return db
}

func TestAuthService_Register_Success(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, token, err := svc.Register(req)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)
	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, user.ID.String(), claims["sub"])
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	req := dto.RegisterRequest{
		Username: "user1",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, _, err := svc.Register(req)
	require.NoError(t, err)

	req.Username = "user2"
	_, _, err = svc.Register(req)

	assert.Error(t, err)
}

func TestAuthService_Login_Success(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	req := dto.RegisterRequest{
		Username: "user",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, _, err := svc.Register(req)
	require.NoError(t, err)

	user, token, err := svc.Login(dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "user", user.Username)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	_, _, err := svc.Login(dto.LoginRequest{
		Email:    "invalid@email.com",
		Password: "password123",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	req := dto.RegisterRequest{
		Username: "user",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, _, err := svc.Register(req)
	require.NoError(t, err)

	_, _, err = svc.Login(dto.LoginRequest{
		Email:    "invalid@email.com",
		Password: "invalidPassword",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestAuthService_GenerateToken(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, _, err := svc.Register(req)

	token, err := svc.GenerateToken(user)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	require.NoError(t, err)

	claims := parsedToken.Claims.(jwt.MapClaims)
	assert.Equal(t, user.ID.String(), claims["sub"])
	assert.Contains(t, claims, "exp")
	assert.Contains(t, claims, "iat")
}

func TestAuthService_GetUserByID_Success(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	req := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	user, _, err := svc.Register(req)
	require.NoError(t, err)

	foundUser, err := svc.GetUserByID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	db := setupTestAuthDB(t)
	svc := &AuthService{
		DB:        db,
		JWTSecret: "test-secret",
		TokenTTL:  time.Hour,
	}

	_, err := svc.GetUserByID(uuid.New())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
