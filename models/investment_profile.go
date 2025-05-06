package models

type InvestmentProfile struct {
	UserID                               string  `gorm:"primaryKey;type:uuid;not null;unique;index" json:"user_id" db:"user_id"`
	User                                 User    `gorm:"foreignKey:UserID;references:ID;onUpdate:CASCADE;onDelete:CASCADE" json:"user"`
	Age                                  int     `gorm:"nullable" json:"age" db:"age"`
	MaxAcceptableShortTermLossPercentage int     `gorm:"nullable" json:"maxAcceptableShortTermLossPercentage" db:"max_acceptable_short_term_loss_percentage"`
	ExpectedAnnualizedRateOfReturn       int     `gorm:"nullable" json:"expectedAnnualizedRateOfReturn" db:"expected_annualized_rate_of_return"`
	TimeHorizon                          string  `gorm:"nullable" json:"timeHorizon" db:"time_horizon"`
	YearsInvesting                       int     `gorm:"nullable" json:"yearsInvesting" db:"years_investing"`
	MonthlyCashFlow                      float64 `gorm:"nullable" json:"monthlyCashFlow" db:"monthly_cash_flow"`
	DefaultCurrency                      string  `gorm:"nullable" json:"defaultCurrency" db:"default_currency"`
}

func (InvestmentProfile) TableName() string {
	return "investment_profiles"
}

type InvestmentProfileCreateRequest struct {
	Age                                  int     `json:"age"`
	MaxAcceptableShortTermLossPercentage int     `json:"maxAcceptableShortTermLossPercentage"`
	ExpectedAnnualizedRateOfReturn       int     `json:"expectedAnnualizedRateOfReturn"`
	TimeHorizon                          string  `json:"timeHorizon"`
	YearsInvesting                       int     `json:"yearsInvesting"`
	MonthlyCashFlow                      float64 `json:"monthlyCashFlow"`
	DefaultCurrency                      string  `json:"defaultCurrency"`
}
