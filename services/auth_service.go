package services

import (
	"asset-dairy/models"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("Invalid refresh token")
)

type AuthServiceInterface interface {
	SignUp(req *models.SignUpRequest) (*models.User, error)
	SignIn(email, password string) (*models.AuthResponse, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (s *AuthService) SignUp(req *models.SignUpRequest) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	id := uuid.New().String()
	err = s.db.QueryRow(
		`INSERT INTO users (id, email, name, username, password_hash, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		id, req.Email, req.Name, req.Username, string(hashed), time.Now(),
	).Scan(&id)
	if err != nil {
		log.Println("Failed to insert user:", err)
		return nil, errors.New("email may already be registered")
	}

	return &models.User{
		ID:        id,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}, nil
}

// generateAccessToken creates a JWT for a given user ID and email
func generateAccessToken(userID string, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret not set in environment")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // 15 min expiry
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateRefreshToken creates a refresh token JWT for a given user ID and email
func generateRefreshToken(userID string, email string) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return "", errors.New("JWT refresh secret not set in environment")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) SignIn(email, password string) (*models.AuthResponse, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	var id, name, username, passwordHash string
	var createdAt time.Time

	err := s.db.QueryRow(`SELECT id, email, name, username, password_hash, created_at FROM users WHERE email = $1`, email).Scan(&id, &email, &name, &username, &passwordHash, &createdAt)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, refreshToken, err := generateTokens(id, email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User: models.User{
			ID:        id,
			Email:     email,
			Name:      name,
			Username:  username,
			CreatedAt: createdAt,
		},
		RefreshToken: refreshToken,
	}, nil
}

// ValidateRefreshToken checks if the given refresh token is valid
func validateRefreshToken(refreshToken string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateTokens creates new access and refresh tokens for a user
func generateTokens(userID, email string) (string, string, error) {
	accessToken, err := generateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates the existing refresh token and generates new tokens
func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	// First validate the refresh token
	claims, err := validateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Extract user details from claims
	userID, ok1 := claims["user_id"].(string)
	email, ok2 := claims["email"].(string)
	if !ok1 || !ok2 {
		return "", "", ErrInvalidToken
	}

	// Generate new tokens
	return generateTokens(userID, email)
}
