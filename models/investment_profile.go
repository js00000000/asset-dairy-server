package models

type InvestmentProfile struct {
	Id                                   int64  `json:"id" db:"id"`
	UserId                               int64  `json:"user_id" db:"user_id"`
	Age                                  int    `json:"age" db:"age"`
	MaxAcceptableShortTermLossPercentage int    `json:"maxAcceptableShortTermLossPercentage" db:"max_acceptable_short_term_loss_percentage"`
	ExpectedAnnualizedRateOfReturn       int    `json:"expectedAnnualizedRateOfReturn" db:"expected_annualized_rate_of_return"`
	TimeHorizon                          string `json:"timeHorizon" db:"time_horizon"`
	YearsInvesting                       int    `json:"yearsInvesting" db:"years_investing"`
}
