package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"asset-dairy/models"

	"github.com/gin-gonic/gin"
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

	var id int64
	err = h.DB.QueryRow(
		`INSERT INTO users (email, name, username, password_hash, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Email, req.Name, req.Username, string(hashed), time.Now(),
	).Scan(&id)
	if err != nil {
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
