package repositories

import (
	"asset-dairy/models"
	"database/sql"
	"time"

	"github.com/google/uuid"
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
	DB *sql.DB
}

// NewUserRepository 建立 UserRepository 實例
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser 建立新使用者
func (r *UserRepository) CreateUser(req *models.SignUpRequest, hashedPassword string) (string, error) {
	id := uuid.New().String()
	err := r.DB.QueryRow(
		`INSERT INTO users (id, email, name, username, password_hash, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		id, req.Email, req.Name, req.Username, hashedPassword, time.Now(),
	).Scan(&id)
	return id, err
}

// FindByEmail 根據 Email 查詢使用者及密碼雜湊
func (r *UserRepository) FindByEmail(email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string
	err := r.DB.QueryRow(`SELECT id, email, name, username, password_hash, created_at FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.Username, &passwordHash, &user.CreatedAt)
	if err != nil { // 包括 sql.ErrNoRows
		return nil, "", err
	}
	return &user, passwordHash, nil
}

// FindByID 根據 UserID 查詢使用者基本資料
func (r *UserRepository) FindByID(userID string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, name, email, username FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Username)
	if err != nil { // 包括 sql.ErrNoRows
		return nil, err
	}
	return &user, nil
}

// GetPasswordHash 根據 UserID 取得密碼雜湊
func (r *UserRepository) GetPasswordHash(userID string) (string, error) {
	var passwordHash string
	err := r.DB.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&passwordHash)
	return passwordHash, err // 包括 sql.ErrNoRows
}

// UpdatePasswordHash 更新使用者密碼雜湊
func (r *UserRepository) UpdatePasswordHash(userID, hashedPassword string) error {
	_, err := r.DB.Exec("UPDATE users SET password_hash = $1 WHERE id = $2", hashedPassword, userID)
	return err
}

// UpdateUser 更新使用者姓名和用戶名
func (r *UserRepository) UpdateUser(userID string, req *models.UserUpdateRequest) error {
	_, err := r.DB.Exec(`UPDATE users SET name = $1, username = $2 WHERE id = $3`,
		req.Name,
		req.Username,
		userID,
	)
	return err
}

// FindInvestmentProfileByUserID 根據 UserID 查詢投資檔案
func (r *UserRepository) FindInvestmentProfileByUserID(userID string) (*models.InvestmentProfile, error) {
	var profile models.InvestmentProfile
	err := r.DB.QueryRow("SELECT id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing, monthly_cash_flow, default_currency FROM investment_profiles WHERE user_id = $1", userID).
		Scan(&profile.Id, &profile.UserId, &profile.Age, &profile.MaxAcceptableShortTermLossPercentage, &profile.ExpectedAnnualizedRateOfReturn, &profile.TimeHorizon, &profile.YearsInvesting, &profile.MonthlyCashFlow, &profile.DefaultCurrency)
	// 不需要特別處理 sql.ErrNoRows，讓呼叫者處理
	return &profile, err
}

// UpsertInvestmentProfile 新增或更新投資檔案
func (r *UserRepository) UpsertInvestmentProfile(userID string, profile *models.InvestmentProfile) error {
	// 檢查 profile 是否為 nil
	if profile == nil {
		// 可以選擇返回錯誤或執行刪除/清空操作，或忽略
		// 這裡選擇忽略，如果需要清空，需要額外的邏輯
		return nil // 或者返回一個錯誤 indicate nil profile
	}
	// 如果 profile.Id 是空字串或需要生成新的 UUID
	profileID := profile.Id
	if profileID == "" {
		profileID = uuid.New().String()
	}

	_, err := r.DB.Exec(`INSERT INTO investment_profiles (id, user_id, age, max_acceptable_short_term_loss_percentage, expected_annualized_rate_of_return, time_horizon, years_investing, monthly_cash_flow, default_currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			age = EXCLUDED.age,
			max_acceptable_short_term_loss_percentage = EXCLUDED.max_acceptable_short_term_loss_percentage,
			expected_annualized_rate_of_return = EXCLUDED.expected_annualized_rate_of_return,
			time_horizon = EXCLUDED.time_horizon,
			years_investing = EXCLUDED.years_investing,
			monthly_cash_flow = EXCLUDED.monthly_cash_flow,
			default_currency = EXCLUDED.default_currency`,
		profileID, // 使用生成的或提供的 ID
		userID,
		profile.Age,
		profile.MaxAcceptableShortTermLossPercentage,
		profile.ExpectedAnnualizedRateOfReturn,
		profile.TimeHorizon,
		profile.YearsInvesting,
		profile.MonthlyCashFlow,
		profile.DefaultCurrency,
	)
	return err
}
