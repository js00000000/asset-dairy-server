package services

import (
	"asset-dairy/models"
	"database/sql"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ProfileServiceInterface interface {
	GetProfile(userID string) (*models.ProfileResponse, error)
	ChangePassword(userID string, currentPassword, newPassword string) error
	UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error)
}

type ProfileService struct {
	db *sql.DB
}

func NewProfileService(db *sql.DB) *ProfileService {
	return &ProfileService{
		db: db,
	}
}

func (s *ProfileService) GetProfile(userID string) (*models.ProfileResponse, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	// Fetch investment profile
	var investmentProfile models.InvestmentProfile
	err = s.db.QueryRow("SELECT id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing, monthly_cash_flow, default_currency FROM investment_profiles WHERE user_id = $1", userID).
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

func (s *ProfileService) ChangePassword(userID string, currentPassword, newPassword string) error {
	// Fetch current password hash
	var passwordHash string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
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
	_, err = s.db.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", string(hashed), userID)
	return err
}

func (s *ProfileService) UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error) {
	// Update name and username in users table
	_, err := s.db.Exec(`UPDATE users SET name = $1, username = $2 WHERE id = $3`,
		req.Name,
		req.Username,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// Upsert investment profile in investment_profiles table
	_, err = s.db.Exec(`INSERT INTO investment_profiles (id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing, monthly_cash_flow, default_currency)
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
	err = s.db.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
