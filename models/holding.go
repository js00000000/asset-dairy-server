package models

// Holding represents a user's current asset holdings
type Holding struct {
	Ticker       string  `json:"ticker" db:"ticker"`
	Quantity     float64 `json:"quantity" db:"quantity"`
	AveragePrice float64 `json:"averagePrice" db:"average_price"`
	AssetType    string  `json:"assetType" db:"asset_type"`
	Currency     string  `json:"currency" db:"currency"`
}
