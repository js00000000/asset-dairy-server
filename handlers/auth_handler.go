package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"asset-dairy/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	id := uuid.New().String()
	err = h.DB.QueryRow(
		`INSERT INTO users (id, email, name, username, password_hash, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		id, req.Email, req.Name, req.Username, string(hashed), time.Now(),
	).Scan(&id)
	if err != nil {
		log.Println("Failed to insert user:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email may already be registered"})
		return
	}

	user := models.User{
		ID:        id,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
	}
	c.JSON(http.StatusCreated, user)
}

// SignInRequest represents the payload for sign-in
// Accepts email and password only
// Example: {"email": "", "password": ""}
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// generateJWT creates a JWT for a given user ID and email
func generateJWT(userID string, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	log.Println(secret)
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

// SignIn authenticates a user and returns a JWT token and user object
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id, email, name, username, passwordHash string
	var createdAt time.Time

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	err := h.DB.QueryRow(`SELECT id, email, name, username, password_hash, created_at FROM users WHERE email = $1`, req.Email).Scan(&id, &email, &name, &username, &passwordHash, &createdAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateJWT(id, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	refreshToken, err := generateRefreshToken(id, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	// Set refresh token as HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": models.User{
			ID:        id,
			Email:     email,
			Name:      name,
			Username:  username,
			CreatedAt: createdAt,
		},
	})
}

// RefreshToken issues a new access token if the refresh token is valid
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}
	secret := os.Getenv("JWT_REFRESH_SECRET")
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token claims"})
		return
	}
	userID, ok1 := claims["user_id"].(string)
	email, ok2 := claims["email"].(string)
	if !ok1 || !ok2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token data"})
		return
	}
	accessToken, err := generateJWT(userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}
	// Generate new refresh token and set it in the cookie
	newRefreshToken, err := generateRefreshToken(userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}
	c.SetCookie("refresh_token", newRefreshToken, 7*24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

// Logout clears the refresh token cookie
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
