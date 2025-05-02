package handlers

import (
	"net/http"
	"time"

	"asset-dairy/models"
	"asset-dairy/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TradeHandler struct {
	service services.TradeServiceInterface
}

func NewTradeHandler(tradeService services.TradeServiceInterface) *TradeHandler {
	return &TradeHandler{
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
	okAcc, err := h.service.IsAccountOwnedByUser(req.AccountID, userID.(string))
	if err != nil || !okAcc {
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
	if err := h.service.CreateTrade(userID.(string), trade); err != nil {
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
	// Check ownership of trade
	okTrade, err := h.service.IsTradeOwnedByUser(id, userID.(string))
	if err != nil || !okTrade {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found or unauthorized"})
		return
	}
	var req models.TradeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.AccountID != "" {
		okAcc, err := h.service.IsAccountOwnedByUser(req.AccountID, userID.(string))
		if err != nil || !okAcc {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or unauthorized account_id"})
			return
		}
	}
	updatedTrade, err := h.service.UpdateTrade(userID.(string), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trade"})
		return
	}
	if updatedTrade == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}
	c.JSON(http.StatusOK, updatedTrade)
}

// Delete a trade
func (h *TradeHandler) DeleteTrade(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")
	deleted, err := h.service.DeleteTrade(userID.(string), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trade"})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade not found or unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "deleted": true})
}
