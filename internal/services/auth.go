package services

import (
	"errors"
	"fmt"
	"time"

	"uptime-checker/internal/dto"
	"uptime-checker/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB        *gorm.DB
	JWTSecret string
}

func (s *AuthService) Register(req dto.RegisterRequest) (*models.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	var user *models.User

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		var defaultRole models.Role
		if err := tx.Where("name = ?", models.RoleUser).First(&defaultRole).Error; err != nil {
			return fmt.Errorf("default role not found: %w", err)
		}

		user = &models.User{
			Username: req.Username,
			Email:    req.Email,
			Credential: models.Credential{
				PasswordHash: string(hashedPassword),
			},
			Roles: []models.Role{defaultRole},
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*models.User, string, error) {
	var user models.User
	err := s.DB.Preload("Credential").Preload("Roles").Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid email or password")
		}
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Credential.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", errors.New("invalid password")
	}

	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) GetUserByID(id interface{}) (*models.User, error) {
	var user models.User

	err := s.DB.Preload("Roles").Where("id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
