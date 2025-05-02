package repositories

import (
	"database/sql"
	"log"

	"asset-dairy/models"
)

type AuthRepositoryInterface interface {
	CreateUser(user *models.User, passwordHash string) error
	FindUserByEmail(email string) (*models.User, string, error)
}

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) CreateUser(user *models.User, passwordHash string) error {
	err := r.DB.QueryRow(
		`INSERT INTO users (id, email, name, username, password_hash, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		user.ID, user.Email, user.Name, user.Username, passwordHash, user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		log.Println("Failed to insert user:", err)
		return err
	}
	return nil
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string
	err := r.DB.QueryRow(
		`SELECT id, email, name, username, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Username, &passwordHash, &user.CreatedAt)
	if err != nil {
		log.Println("Failed to find user:", err)
		return nil, "", err
	}
	return &user, passwordHash, nil
}
