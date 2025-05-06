package handlers

import (
	"asset-dairy/models"
	"asset-dairy/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.UserSignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.SignUp(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req models.UserSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.SignIn(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token as HttpOnly cookie
	c.SetCookie("refresh_token", response.RefreshToken, 7*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": response.Token,
		"user":  response.User,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	// Call service to validate and generate new tokens
	accessToken, newRefreshToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		if err == services.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new tokens"})
		}
		return
	}

	// Set new refresh token in cookie
	c.SetCookie("refresh_token", newRefreshToken, 7*24*3600, "/", "", false, true)

	// Return new access token
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
