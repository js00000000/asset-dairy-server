package models

type ProfileResponse struct {
	ID                string             `json:"id"`
	Email             string             `json:"email"`
	Name              string             `json:"name"`
	Username          string             `json:"username"`
	InvestmentProfile *InvestmentProfile `json:"investmentProfile,omitempty"`
}
