package handlers

import (
	"database/sql"
	"net/http"

	"asset-dairy/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	DB *sql.DB
}

func NewProfileHandler(db *sql.DB) *ProfileHandler {
	return &ProfileHandler{DB: db}
}

// GetProfile returns the current user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user models.User
	err := h.DB.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Fetch investment profile
	var investmentProfile models.InvestmentProfile
	err = h.DB.QueryRow("SELECT id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing FROM investment_profiles WHERE user_id = $1", userID).
		Scan(&investmentProfile.Id, &investmentProfile.UserId, &investmentProfile.Age, &investmentProfile.MaxAcceptableShortTermLossPercentage, &investmentProfile.ExpectedAnnualizedRateOfReturn, &investmentProfile.TimeHorizon, &investmentProfile.YearsInvesting)
	if err == sql.ErrNoRows {
		// No investment profile found, return nil for this field
		c.JSON(http.StatusOK, models.ProfileResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Username: user.Username,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch investment profile"})
		return
	}

	c.JSON(http.StatusOK, models.ProfileResponse{
		ID:                user.ID,
		Email:             user.Email,
		Name:              user.Name,
		Username:          user.Username,
		InvestmentProfile: &investmentProfile,
	})
}

// UpdateProfile updates the current user's profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update name and username in users table
	_, err := h.DB.Exec(`UPDATE users SET name = $1, username = $2 WHERE id = $3`,
		req.Name,
		req.Username,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Upsert investment profile in investment_profiles table
	_, err = h.DB.Exec(`INSERT INTO investment_profiles (id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			age = EXCLUDED.age,
			max_acceptable_short_term_loss_percentage = EXCLUDED.max_acceptable_short_term_loss_percentage,
			expected_annualized_rate_of_return = EXCLUDED.expected_annualized_rate_of_return,
			time_horizon = EXCLUDED.time_horizon,
			years_investing = EXCLUDED.years_investing`,
		uuid.New().String(),
		userID,
		req.InvestmentProfile.Age,
		req.InvestmentProfile.MaxAcceptableShortTermLossPercentage,
		req.InvestmentProfile.ExpectedAnnualizedRateOfReturn,
		req.InvestmentProfile.TimeHorizon,
		req.InvestmentProfile.YearsInvesting,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update investment profile"})
		return
	}

	var user models.User
	h.DB.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	c.JSON(http.StatusOK, user)
}
