package services

import (
	"asset-dairy/models"
	"asset-dairy/repositories"
	"errors"
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
	SignUp(req *models.UserSignUpRequest) (*models.User, error)
	SignIn(email, password string) (*models.AuthResponse, error)
	RefreshToken(refreshToken string) (string, string, error)
}

type AuthService struct {
	repo repositories.AuthRepositoryInterface
}

func NewAuthService(repo repositories.AuthRepositoryInterface) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) SignUp(req *models.UserSignUpRequest) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Name:      req.Name,
		Username:  req.Username,
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateUser(user, string(hashed))
	if err != nil {
		return nil, errors.New("email may already be registered")
	}

	return user, nil
}

func (s *AuthService) SignIn(email, password string) (*models.AuthResponse, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, passwordHash, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, refreshToken, err := generateTokens(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:        token,
		User:         *user,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := validateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	userID, ok1 := claims["user_id"].(string)
	email, ok2 := claims["email"].(string)
	if !ok1 || !ok2 {
		return "", "", ErrInvalidToken
	}

	return generateTokens(userID, email)
}

func generateAccessToken(userID string, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret not set in environment")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateRefreshToken(userID string, email string) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		return "", errors.New("JWT refresh secret not set in environment")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

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
