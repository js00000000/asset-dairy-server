package repositories

import (
	"asset-dairy/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryInterface 定義了使用者和個人檔案資料庫操作的介面
type UserRepositoryInterface interface {
	CreateUser(req *models.SignUpRequest, hashedPassword string) (string, error)
	FindByEmail(email string) (*models.User, string, error) // 回傳 User 和 passwordHash
	FindByID(userID string) (*models.User, error)
	GetPasswordHash(userID string) (string, error)
	UpdatePasswordHash(userID, hashedPassword string) error
	UpdateUser(userID string, req *models.UserUpdateRequest) error
	FindInvestmentProfileByUserID(userID string) (*models.InvestmentProfile, error)
	UpsertInvestmentProfile(userID string, profile *models.InvestmentProfile) error
}

// UserRepository 實作了 UserRepositoryInterface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 建立 UserRepository 實例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser 建立新使用者
func (r *UserRepository) CreateUser(req *models.SignUpRequest, hashedPassword string) (string, error) {
	gormUser := &models.GormUser{
		ID:            uuid.New().String(),
		Email:         req.Email,
		Name:          req.Name,
		Username:      req.Username,
		Password_Hash: hashedPassword,
	}
	result := r.db.Create(gormUser)
	return gormUser.ID, result.Error
}

// FindByEmail 根據 Email 查詢使用者及密碼雜湊
func (r *UserRepository) FindByEmail(email string) (*models.User, string, error) {
	var gormUser models.GormUser
	result := r.db.Where(&models.GormUser{Email: email}).First(&gormUser)
	if result.Error != nil {
		return nil, "", result.Error
	}

	user := &models.User{
		ID:       gormUser.ID,
		Email:    gormUser.Email,
		Name:     gormUser.Name,
		Username: gormUser.Username,
		// CreatedAt removed due to model changes
	}

	return user, gormUser.Password_Hash, nil
}

// FindByID 根據 UserID 查詢使用者基本資料
func (r *UserRepository) FindByID(userID string) (*models.User, error) {
	var gormUser models.GormUser
	result := r.db.Select("id", "name", "email", "username").Where("id = ?", userID).First(&gormUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &models.User{
		ID:       gormUser.ID,
		Name:     gormUser.Name,
		Email:    gormUser.Email,
		Username: gormUser.Username,
	}, nil
}

// GetPasswordHash 根據 UserID 取得密碼雜湊
func (r *UserRepository) GetPasswordHash(userID string) (string, error) {
	var gormUser models.GormUser
	result := r.db.Select("password").Where(&models.GormUser{ID: userID}).First(&gormUser)
	return gormUser.Password_Hash, result.Error
}

// UpdatePasswordHash 更新使用者密碼雜湊
func (r *UserRepository) UpdatePasswordHash(userID, hashedPassword string) error {
	result := r.db.Model(&models.GormUser{}).Where(&models.GormUser{ID: userID}).Update("password", hashedPassword)
	return result.Error
}

// UpdateUser 更新使用者姓名和用戶名
func (r *UserRepository) UpdateUser(userID string, req *models.UserUpdateRequest) error {
	result := r.db.Model(&models.GormUser{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"name":     req.Name,
		"username": req.Username,
	})
	return result.Error
}

// FindInvestmentProfileByUserID 根據 UserID 查詢投資檔案
func (r *UserRepository) FindInvestmentProfileByUserID(userID string) (*models.InvestmentProfile, error) {
	var gormProfile models.GormInvestmentProfile
	result := r.db.Where(&models.GormInvestmentProfile{UserID: userID}).First(&gormProfile)
	if result.Error != nil {
		return nil, result.Error
	}

	return &models.InvestmentProfile{
		Id:                                   gormProfile.ID,
		UserId:                               gormProfile.UserID,
		Age:                                  gormProfile.Age,
		MaxAcceptableShortTermLossPercentage: int(gormProfile.MaxAcceptableShortTermLossPercentage),
		ExpectedAnnualizedRateOfReturn:       int(gormProfile.ExpectedAnnualizedRateOfReturn),
		TimeHorizon:                          gormProfile.TimeHorizon,
		YearsInvesting:                       gormProfile.YearsInvesting,
		MonthlyCashFlow:                      gormProfile.MonthlyCashFlow,
		DefaultCurrency:                      gormProfile.DefaultCurrency,
	}, nil
}

// UpsertInvestmentProfile 新增或更新投資檔案
func (r *UserRepository) UpsertInvestmentProfile(userID string, profile *models.InvestmentProfile) error {
	if profile == nil {
		return nil
	}

	gormProfile := &models.GormInvestmentProfile{
		UserID:                               userID,
		Age:                                  profile.Age,
		MaxAcceptableShortTermLossPercentage: int(profile.MaxAcceptableShortTermLossPercentage),
		ExpectedAnnualizedRateOfReturn:       int(profile.ExpectedAnnualizedRateOfReturn),
		TimeHorizon:                          profile.TimeHorizon,
		YearsInvesting:                       profile.YearsInvesting,
		MonthlyCashFlow:                      profile.MonthlyCashFlow,
		DefaultCurrency:                      profile.DefaultCurrency,
	}

	// Use FirstOrCreate to either find an existing record or create a new one
	result := r.db.Where(&models.GormInvestmentProfile{UserID: userID}).Assign(gormProfile).FirstOrCreate(gormProfile)
	return result.Error
}
