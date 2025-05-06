package repositories

import (
	"log"
	"time"

	"asset-dairy/models"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	CreateUser(user *models.User, passwordHash string) error
	FindUserByEmail(email string) (*models.User, string, error)
}

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) CreateUser(user *models.User, passwordHash string) error {
	gormUser := &models.User{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		Username:      user.Username,
		Password_Hash: passwordHash,
	}

	result := r.DB.Create(gormUser)
	if result.Error != nil {
		log.Println("Failed to insert user:", result.Error)
		return result.Error
	}
	return nil
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, string, error) {
	var gormUser models.User
	result := r.DB.Where(&models.User{Email: email}).First(&gormUser)
	if result.Error != nil {
		log.Println("Failed to find user:", result.Error)
		return nil, "", result.Error
	}

	user := &models.User{
		ID:        gormUser.ID,
		Email:     gormUser.Email,
		Name:      gormUser.Name,
		Username:  gormUser.Username,
		CreatedAt: time.Now(), // Use current time as a fallback
	}

	return user, gormUser.Password_Hash, nil
}
