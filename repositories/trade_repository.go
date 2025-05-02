package repositories

import (
	"asset-dairy/models"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
)

// TradeRepositoryInterface defines methods for trade-related database operations
type TradeRepositoryInterface interface {
	ListTrades(userID string) ([]models.Trade, error)
	CreateTrade(userID string, trade models.Trade) error
	UpdateTrade(userID, tradeID string, req models.TradeUpdateRequest) (*models.Trade, error)
	DeleteTrade(userID, tradeID string) (bool, error)
	IsAccountOwnedByUser(accountID, userID string) (bool, error)
	IsTradeOwnedByUser(tradeID, userID string) (bool, error)
}

// TradeRepository implements TradeRepositoryInterface
type TradeRepository struct {
	db *sql.DB
}

// NewTradeRepository creates a new TradeRepository instance
func NewTradeRepository(db *sql.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// ListTrades retrieves all trades for a given user
func (r *TradeRepository) ListTrades(userID string) ([]models.Trade, error) {
	rows, err := r.db.Query("SELECT id, type, asset_type, ticker, trade_date, quantity, price, currency, account_id, reason FROM trades WHERE user_id = $1", userID)
	if err != nil {
		log.Println("TradeRepository: Failed to fetch trades:", err)
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

func (r *TradeRepository) IsAccountOwnedByUser(accountID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM accounts WHERE id = $1 AND user_id = $2`, accountID, userID).Scan(&count)
	return count > 0, err
}

func (r *TradeRepository) IsTradeOwnedByUser(tradeID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM trades t JOIN accounts a ON t.account_id = a.id WHERE t.id = $1 AND a.user_id = $2`, tradeID, userID).Scan(&count)
	return count > 0, err
}

func (r *TradeRepository) CreateTrade(userID string, trade models.Trade) error {
	_, err := r.db.Exec(`INSERT INTO trades (id, user_id, type, asset_type, ticker, trade_date, quantity, price, currency, account_id, reason) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		trade.ID, userID, trade.Type, trade.AssetType, trade.Ticker, trade.TradeDate, trade.Quantity, trade.Price, trade.Currency, trade.AccountID, trade.Reason,
	)
	return err
}

func (r *TradeRepository) UpdateTrade(userID, tradeID string, req models.TradeUpdateRequest) (*models.Trade, error) {
	// Build update query dynamically based on non-zero fields
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1
	if req.Type != "" {
		setClauses = append(setClauses, "type = $"+strconv.Itoa(argIdx))
		args = append(args, req.Type)
		argIdx++
	}
	if req.AssetType != "" {
		setClauses = append(setClauses, "asset_type = $"+strconv.Itoa(argIdx))
		args = append(args, req.AssetType)
		argIdx++
	}
	if req.Ticker != "" {
		setClauses = append(setClauses, "ticker = $"+strconv.Itoa(argIdx))
		args = append(args, req.Ticker)
		argIdx++
	}
	if req.TradeDate != "" {
		tradeDate, err := time.Parse("2006-01-02", req.TradeDate)
		if err != nil {
			return nil, err
		}
		setClauses = append(setClauses, "trade_date = $"+strconv.Itoa(argIdx))
		args = append(args, tradeDate)
		argIdx++
	}
	if req.Quantity != 0 {
		setClauses = append(setClauses, "quantity = $"+strconv.Itoa(argIdx))
		args = append(args, req.Quantity)
		argIdx++
	}
	if req.Price != 0 {
		setClauses = append(setClauses, "price = $"+strconv.Itoa(argIdx))
		args = append(args, req.Price)
		argIdx++
	}
	if req.Currency != "" {
		setClauses = append(setClauses, "currency = $"+strconv.Itoa(argIdx))
		args = append(args, req.Currency)
		argIdx++
	}
	if req.AccountID != "" {
		setClauses = append(setClauses, "account_id = $"+strconv.Itoa(argIdx))
		args = append(args, req.AccountID)
		argIdx++
	}
	if req.Reason != nil {
		setClauses = append(setClauses, "reason = $"+strconv.Itoa(argIdx))
		args = append(args, req.Reason)
		argIdx++
	}
	if len(setClauses) == 0 {
		return nil, nil // No fields to update
	}
	args = append(args, tradeID)
	query := "UPDATE trades SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(argIdx)
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	// Return updated trade
	var t models.Trade
	var reason sql.NullString
	err = r.db.QueryRow(`SELECT t.id, t.type, t.asset_type, t.ticker, t.trade_date, t.quantity, t.price, t.currency, t.account_id, t.reason
		FROM trades t JOIN accounts a ON t.account_id = a.id WHERE t.id = $1 AND a.user_id = $2`, tradeID, userID).Scan(
		&t.ID, &t.Type, &t.AssetType, &t.Ticker, &t.TradeDate, &t.Quantity, &t.Price, &t.Currency, &t.AccountID, &reason,
	)
	if err != nil {
		return nil, err
	}
	if reason.Valid {
		t.Reason = &reason.String
	}
	return &t, nil
}

func (r *TradeRepository) DeleteTrade(userID, tradeID string) (bool, error) {
	res, err := r.db.Exec(`DELETE FROM trades WHERE id = $1 AND account_id IN (SELECT id FROM accounts WHERE user_id = $2)`, tradeID, userID)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
