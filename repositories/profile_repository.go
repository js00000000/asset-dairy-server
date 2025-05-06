package repositories

import (
	"log"

	"asset-dairy/models"

	"gorm.io/gorm"

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
	db *gorm.DB
}

// NewProfileRepository creates a new ProfileRepository instance
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetProfile retrieves user profile and investment profile from the database
func (r *ProfileRepository) GetProfile(userID string) (*models.ProfileResponse, error) {
	var gormUser models.GormUser
	result := r.db.Where(&models.GormUser{ID: userID}).First(&gormUser)
	if result.Error != nil {
		log.Println("Failed to fetch user:", result.Error)
		return nil, result.Error
	}

	// Fetch investment profile
	var gormInvestmentProfile models.GormInvestmentProfile
	result = r.db.Where(&models.GormInvestmentProfile{UserID: userID}).First(&gormInvestmentProfile)
	if result.Error == gorm.ErrRecordNotFound {
		return &models.ProfileResponse{
			ID:       gormUser.ID,
			Email:    gormUser.Email,
			Name:     gormUser.Name,
			Username: gormUser.Username,
		}, nil
	} else if result.Error != nil {
		log.Println("Failed to fetch investment profile:", result.Error)
		return nil, result.Error
	}

	return &models.ProfileResponse{
		ID:       gormUser.ID,
		Email:    gormUser.Email,
		Name:     gormUser.Name,
		Username: gormUser.Username,
		InvestmentProfile: &models.InvestmentProfile{
			Id:                                   gormInvestmentProfile.ID,
			UserId:                               gormInvestmentProfile.UserID,
			Age:                                  int(gormInvestmentProfile.Age),
			MaxAcceptableShortTermLossPercentage: int(gormInvestmentProfile.MaxAcceptableShortTermLossPercentage),
			ExpectedAnnualizedRateOfReturn:       int(gormInvestmentProfile.ExpectedAnnualizedRateOfReturn),
			TimeHorizon:                          gormInvestmentProfile.TimeHorizon,
			YearsInvesting:                       int(gormInvestmentProfile.YearsInvesting),
			MonthlyCashFlow:                      gormInvestmentProfile.MonthlyCashFlow,
			DefaultCurrency:                      gormInvestmentProfile.DefaultCurrency,
		},
	}, nil
}

// ChangePassword updates the user's password after verifying the current password
func (r *ProfileRepository) ChangePassword(userID string, currentPassword, newPassword string) error {
	var gormUser models.GormUser
	result := r.db.Where(&models.GormUser{ID: userID}).First(&gormUser)
	if result.Error != nil {
		log.Println("Failed to find user:", result.Error)
		return result.Error
	}

	// Check current password
	err := bcrypt.CompareHashAndPassword([]byte(gormUser.Password_Hash), []byte(currentPassword))
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	gormUser.Password_Hash = string(hashedPassword)
	result = r.db.Save(&gormUser)
	if result.Error != nil {
		log.Println("Failed to update password:", result.Error)
		return result.Error
	}

	return nil
}

func (r *ProfileRepository) UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error) {
	// Find and update user
	var gormUser models.GormUser
	result := r.db.Where(&models.GormUser{ID: userID}).First(&gormUser)
	if result.Error != nil {
		log.Println("Failed to find user:", result.Error)
		return nil, result.Error
	}

	// Update fields
	gormUser.Name = req.Name
	gormUser.Username = req.Username

	result = r.db.Save(&gormUser)
	if result.Error != nil {
		log.Println("Failed to update user:", result.Error)
		return nil, result.Error
	}

	// Upsert investment profile
	if req.InvestmentProfile != nil {
		// Check if investment profile already exists
		var existingProfile models.GormInvestmentProfile
		result = r.db.Where(&models.GormInvestmentProfile{UserID: userID}).First(&existingProfile)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new investment profile
			newProfile := models.GormInvestmentProfile{
				ID:                                   uuid.New().String(),
				UserID:                               userID,
				Age:                                  int(req.InvestmentProfile.Age),
				MaxAcceptableShortTermLossPercentage: float64(req.InvestmentProfile.MaxAcceptableShortTermLossPercentage),
				ExpectedAnnualizedRateOfReturn:       float64(req.InvestmentProfile.ExpectedAnnualizedRateOfReturn),
				TimeHorizon:                          req.InvestmentProfile.TimeHorizon,
				YearsInvesting:                       int(req.InvestmentProfile.YearsInvesting),
				MonthlyCashFlow:                      req.InvestmentProfile.MonthlyCashFlow,
				DefaultCurrency:                      req.InvestmentProfile.DefaultCurrency,
			}
			result = r.db.Create(&newProfile)
		} else if result.Error == nil {
			// Update existing investment profile
			existingProfile.Age = int(req.InvestmentProfile.Age)
			existingProfile.MaxAcceptableShortTermLossPercentage = float64(req.InvestmentProfile.MaxAcceptableShortTermLossPercentage)
			existingProfile.ExpectedAnnualizedRateOfReturn = float64(req.InvestmentProfile.ExpectedAnnualizedRateOfReturn)
			existingProfile.TimeHorizon = req.InvestmentProfile.TimeHorizon
			existingProfile.YearsInvesting = int(req.InvestmentProfile.YearsInvesting)
			existingProfile.MonthlyCashFlow = req.InvestmentProfile.MonthlyCashFlow
			existingProfile.DefaultCurrency = req.InvestmentProfile.DefaultCurrency
			result = r.db.Save(&existingProfile)
		} else {
			log.Println("Failed to process investment profile:", result.Error)
			return nil, result.Error
		}

		if result.Error != nil {
			log.Println("Failed to save investment profile:", result.Error)
			return nil, result.Error
		}
	}

	return &models.User{
		ID:       gormUser.ID,
		Name:     gormUser.Name,
		Email:    gormUser.Email,
		Username: gormUser.Username,
	}, nil
}
