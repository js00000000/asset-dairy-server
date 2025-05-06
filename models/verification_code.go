package models

import (
	"time"
)

type VerificationCode struct {
	Email     string    `gorm:"primarykey;unique;index"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
