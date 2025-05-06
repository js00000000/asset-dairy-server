package models

import "time"

type GormAccount struct {
	ID       string  `gorm:"primaryKey;type:uuid"`
	UserID   string  `gorm:"type:uuid;not null;index"`
	Name     string  `gorm:"not null"`
	Currency string  `gorm:"not null"`
	Balance  float64 `gorm:"not null"`
}

func (GormAccount) TableName() string {
	return "accounts"
}

type GormUser struct {
	ID            string `gorm:"primaryKey;type:uuid"`
	Name          string `gorm:"not null"`
	Email         string `gorm:"unique;not null"`
	Username      string `gorm:"unique;not null"`
	Password_Hash string `gorm:"not null"`
}

func (GormUser) TableName() string {
	return "users"
}

type GormInvestmentProfile struct {
	ID                                   string  `gorm:"primaryKey;type:uuid"`
	UserID                               string  `gorm:"type:uuid;not null;index"`
	Age                                  int     `gorm:"not null"`
	MaxAcceptableShortTermLossPercentage float64 `gorm:"not null"`
	ExpectedAnnualizedRateOfReturn       float64 `gorm:"not null"`
	TimeHorizon                          string  `gorm:"not null"`
	YearsInvesting                       int     `gorm:"not null"`
	MonthlyCashFlow                      float64 `gorm:"not null"`
	DefaultCurrency                      string  `gorm:"not null"`
}

func (GormInvestmentProfile) TableName() string {
	return "investment_profiles"
}

type GormTrade struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	Type      string    `gorm:"not null"`
	AssetType string    `gorm:"not null"`
	Ticker    string    `gorm:"not null"`
	TradeDate time.Time `gorm:"not null"`
	Quantity  float64   `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Currency  string    `gorm:"not null"`
	AccountID string    `gorm:"type:uuid;not null;index"`
	Reason    *string   `gorm:"nullable"`
}

func (GormTrade) TableName() string {
	return "trades"
}
