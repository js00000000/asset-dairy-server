package models

import "time"

type GormAccount struct {
	ID       string   `gorm:"primaryKey;type:uuid"`
	UserID   string   `gorm:"type:uuid;not null;index"`
	User     GormUser `gorm:"foreignKey:UserID;references:ID;onUpdate:CASCADE;onDelete:CASCADE"`
	Name     string   `gorm:"not null"`
	Currency string   `gorm:"not null"`
	Balance  float64  `gorm:"not null"`
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
	ID                                   string   `gorm:"primaryKey;type:uuid"`
	UserID                               string   `gorm:"type:uuid;not null;unique;index"`
	User                                 GormUser `gorm:"foreignKey:UserID;references:ID;onUpdate:CASCADE;onDelete:CASCADE"`
	Age                                  int      `gorm:"nullable"`
	MaxAcceptableShortTermLossPercentage int      `gorm:"nullable"`
	ExpectedAnnualizedRateOfReturn       int      `gorm:"nullable"`
	TimeHorizon                          string   `gorm:"nullable"`
	YearsInvesting                       int      `gorm:"nullable"`
	MonthlyCashFlow                      float64  `gorm:"nullable"`
	DefaultCurrency                      string   `gorm:"nullable"`
}

func (GormInvestmentProfile) TableName() string {
	return "investment_profiles"
}

type GormTrade struct {
	ID        string      `gorm:"primaryKey;type:uuid"`
	UserID    string      `gorm:"type:uuid;not null;index"`
	User      GormUser    `gorm:"foreignKey:UserID;references:ID;onUpdate:CASCADE;onDelete:CASCADE"`
	Type      string      `gorm:"not null"`
	AssetType string      `gorm:"not null"`
	Ticker    string      `gorm:"not null"`
	TradeDate time.Time   `gorm:"not null"`
	Quantity  float64     `gorm:"not null"`
	Price     float64     `gorm:"not null"`
	Currency  string      `gorm:"not null"`
	AccountID string      `gorm:"type:uuid;not null;index"`
	Account   GormAccount `gorm:"foreignKey:AccountID;references:ID;onUpdate:CASCADE"`
	Reason    *string     `gorm:"nullable"`
}

func (GormTrade) TableName() string {
	return "trades"
}
