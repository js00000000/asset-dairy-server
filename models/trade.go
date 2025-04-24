package models

import "time"

// Trade represents a trade transaction.
type Trade struct {
	ID        string    `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"` // buy or sell
	AssetType string    `json:"assetType" db:"asset_type"` // stock or crypto
	Ticker    string    `json:"ticker" db:"ticker"`
	TradeDate time.Time `json:"tradeDate" db:"trade_date"`
	Quantity  float64   `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	AccountID string    `json:"accountId" db:"account_id"`
	Reason    *string   `json:"reason,omitempty" db:"reason"`
}

// TradeCreateRequest for creating a trade
// (optional: can be used for binding in handlers)
type TradeCreateRequest struct {
	Type      string    `json:"type" binding:"required,oneof=buy sell"`
	AssetType string    `json:"assetType" binding:"required,oneof=stock crypto"`
	Ticker    string    `json:"ticker" binding:"required"`
	TradeDate string    `json:"tradeDate" binding:"required"`
	Quantity  float64   `json:"quantity" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
	AccountID string    `json:"accountId" binding:"required"`
	Reason    *string   `json:"reason"`
}

type TradeUpdateRequest struct {
	Type      string    `json:"type" binding:"omitempty,oneof=buy sell"`
	AssetType string    `json:"assetType" binding:"omitempty,oneof=stock crypto"`
	Ticker    string    `json:"ticker" binding:"omitempty"`
	TradeDate string    `json:"tradeDate" binding:"omitempty"`
	Quantity  float64   `json:"quantity" binding:"omitempty"`
	Price     float64   `json:"price" binding:"omitempty"`
	AccountID string    `json:"accountId" binding:"omitempty"`
	Reason    *string   `json:"reason"`
}
