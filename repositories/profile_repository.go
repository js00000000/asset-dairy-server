package repositories

import (
	"asset-dairy/models"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ProfileRepositoryInterface defines methods for profile-related database operations
type ProfileRepositoryInterface interface {
	GetProfile(userID string) (*models.ProfileResponse, error)
	ChangePassword(userID string, currentPassword, newPassword string) error
	UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error)
}

// ProfileRepository implements ProfileRepositoryInterface
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository creates a new ProfileRepository instance
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetProfile retrieves user profile and investment profile from the database
func (r *ProfileRepository) GetProfile(userID string) (*models.ProfileResponse, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	// Fetch investment profile
	var investmentProfile models.InvestmentProfile
	err = r.db.QueryRow("SELECT id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing, monthly_cash_flow, default_currency FROM investment_profiles WHERE user_id = $1", userID).
		Scan(&investmentProfile.Id, &investmentProfile.UserId, &investmentProfile.Age, &investmentProfile.MaxAcceptableShortTermLossPercentage, &investmentProfile.ExpectedAnnualizedRateOfReturn, &investmentProfile.TimeHorizon, &investmentProfile.YearsInvesting, &investmentProfile.MonthlyCashFlow, &investmentProfile.DefaultCurrency)
	if err == sql.ErrNoRows {
		return &models.ProfileResponse{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Username: user.Username,
		}, nil
	} else if err != nil {
		return nil, err
	}

	return &models.ProfileResponse{
		ID:                user.ID,
		Email:             user.Email,
		Name:              user.Name,
		Username:          user.Username,
		InvestmentProfile: &investmentProfile,
	}, nil
}

// ChangePassword updates the user's password after verifying the current password
func (r *ProfileRepository) ChangePassword(userID string, currentPassword, newPassword string) error {
	// Fetch current password hash
	var passwordHash string
	err := r.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
	if err != nil {
		return err
	}

	// Compare current password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(currentPassword)); err != nil {
		return err
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password in DB
	_, err = r.db.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", string(hashed), userID)
	return err
}

// UpdateProfile updates user details and investment profile
func (r *ProfileRepository) UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error) {
	// Update name and username in users table
	_, err := r.db.Exec(`UPDATE users SET name = $1, username = $2 WHERE id = $3`,
		req.Name,
		req.Username,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// Upsert investment profile in investment_profiles table
	_, err = r.db.Exec(`INSERT INTO investment_profiles (id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing, monthly_cash_flow, default_currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			age = EXCLUDED.age,
			max_acceptable_short_term_loss_percentage = EXCLUDED.max_acceptable_short_term_loss_percentage,
			expected_annualized_rate_of_return = EXCLUDED.expected_annualized_rate_of_return,
			time_horizon = EXCLUDED.time_horizon,
			years_investing = EXCLUDED.years_investing,
			monthly_cash_flow = EXCLUDED.monthly_cash_flow,
			default_currency = EXCLUDED.default_currency`,
		uuid.New().String(),
		userID,
		req.InvestmentProfile.Age,
		req.InvestmentProfile.MaxAcceptableShortTermLossPercentage,
		req.InvestmentProfile.ExpectedAnnualizedRateOfReturn,
		req.InvestmentProfile.TimeHorizon,
		req.InvestmentProfile.YearsInvesting,
		req.InvestmentProfile.MonthlyCashFlow,
		req.InvestmentProfile.DefaultCurrency,
	)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.db.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
