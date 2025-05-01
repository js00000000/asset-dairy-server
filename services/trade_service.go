package services

import (
	"asset-dairy/models"
	"database/sql"
	"log"
)

type TradeServiceInterface interface {
	ListTrades(userID string) ([]models.Trade, error)
}

type TradeService struct {
	db *sql.DB
}

func NewTradeService(db *sql.DB) *TradeService {
	return &TradeService{db: db}
}

// ListTrades retrieves all trades for a given user
func (s *TradeService) ListTrades(userID string) ([]models.Trade, error) {
	rows, err := s.db.Query("SELECT id, type, asset_type, ticker, trade_date, quantity, price, currency, account_id, reason FROM trades WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Failed to fetch trades:", err)
		return nil, err
	}
	defer rows.Close()

	trades := []models.Trade{}
	for rows.Next() {
		var t models.Trade
		var reason sql.NullString
		if err := rows.Scan(&t.ID, &t.Type, &t.AssetType, &t.Ticker, &t.TradeDate, &t.Quantity, &t.Price, &t.Currency, &t.AccountID, &reason); err != nil {
			return nil, err
		}
		if reason.Valid {
			t.Reason = &reason.String
		}
		trades = append(trades, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}
