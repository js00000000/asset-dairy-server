package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"asset-dairy/models"
	"asset-dairy/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TradeHandler struct {
	DB      *sql.DB
	service services.TradeServiceInterface
}

func NewTradeHandler(db *sql.DB, tradeService services.TradeServiceInterface) *TradeHandler {
	return &TradeHandler{
		DB:      db,
		service: tradeService,
	}
}

// List all trades for a given account or user
func (h *TradeHandler) ListTrades(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trades, err := h.service.ListTrades(userID.(string))
	if err != nil {
		log.Println("Failed to fetch trades:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// Create a trade
func (h *TradeHandler) CreateTrade(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req models.TradeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Verify account belongs to user
	var count int
	err := h.DB.QueryRow(`SELECT COUNT(*) FROM accounts WHERE id = $1 AND user_id = $2`, req.AccountID, userID).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or unauthorized account_id"})
		return
	}
	tradeDate, err := time.Parse("2006-01-02", req.TradeDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tradeDate format, use YYYY-MM-DD"})
		return
	}
	trade := models.Trade{
		ID:        uuid.New().String(),
		Type:      req.Type,
		AssetType: req.AssetType,
		Ticker:    req.Ticker,
		TradeDate: tradeDate,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Currency:  req.Currency,
		AccountID: req.AccountID,
		Reason:    req.Reason,
	}
	_, err = h.DB.Exec(`INSERT INTO trades (id, user_id, type, asset_type, ticker, trade_date, quantity, price, currency, account_id, reason) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		trade.ID, userID, trade.Type, trade.AssetType, trade.Ticker, trade.TradeDate, trade.Quantity, trade.Price, trade.Currency, trade.AccountID, trade.Reason,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trade"})
		return
	}
	c.JSON(http.StatusCreated, trade)
}

// Update a trade
func (h *TradeHandler) UpdateTrade(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	// Check ownership of trade by joining through account
	var count int
	err := h.DB.QueryRow(`SELECT COUNT(*) FROM trades t JOIN accounts a ON t.account_id = a.id WHERE t.id = $1 AND a.user_id = $2`, id, userID).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found or unauthorized"})
		return
	}
	var req models.TradeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tradeDate format, use YYYY-MM-DD"})
			return
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
		// Verify new account belongs to user
		var accCount int
		err := h.DB.QueryRow(`SELECT COUNT(*) FROM accounts WHERE id = $1 AND user_id = $2`, req.AccountID, userID).Scan(&accCount)
		if err != nil || accCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or unauthorized account_id"})
			return
		}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}
	args = append(args, id)
	query := "UPDATE trades SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(argIdx)
	_, err = h.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trade"})
		return
	}
	// Return updated trade
	var t models.Trade
	var reason sql.NullString
	err = h.DB.QueryRow(`SELECT t.id, t.type, t.asset_type, t.ticker, t.trade_date, t.quantity, t.price, t.currency, t.account_id, t.reason
		FROM trades t JOIN accounts a ON t.account_id = a.id WHERE t.id = $1 AND a.user_id = $2`, id, userID).Scan(
		&t.ID, &t.Type, &t.AssetType, &t.Ticker, &t.TradeDate, &t.Quantity, &t.Price, &t.Currency, &t.AccountID, &reason,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated trade"})
		return
	}
	if reason.Valid {
		t.Reason = &reason.String
	}
	c.JSON(http.StatusOK, t)
}

// Delete a trade
func (h *TradeHandler) DeleteTrade(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	// Only delete if trade belongs to user (via account)
	res, err := h.DB.Exec(`DELETE FROM trades WHERE id = $1 AND account_id IN (SELECT id FROM accounts WHERE user_id = $2)`, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trade"})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found or unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "deleted": true})
}
