package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"asset-dairy/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	DB *sql.DB
}

func NewAccountHandler(db *sql.DB) *AccountHandler {
	return &AccountHandler{DB: db}
}

// ListAccounts returns all accounts for the current user
func (h *AccountHandler) ListAccounts(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	rows, err := h.DB.Query("SELECT id, name, currency, balance FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Failed to fetch accounts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}
	defer rows.Close()
	accounts := []models.Account{}
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Balance); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan account"})
			return
		}
		accounts = append(accounts, acc)
	}
	c.JSON(http.StatusOK, accounts)
}

// CreateAccount creates a new account for the current user
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req models.AccountCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id string
	err := h.DB.QueryRow("INSERT INTO accounts (id, user_id, name, currency, balance) VALUES ($1, $2, $3, $4, $5) RETURNING id", uuid.New().String(), userID, req.Name, req.Currency, req.Balance).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}
	acc := models.Account{ID: id, Name: req.Name, Currency: req.Currency, Balance: req.Balance}
	c.JSON(http.StatusCreated, acc)
}

// UpdateAccount updates an account by id for the current user
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	accID := c.Param("id")
	var req models.AccountUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.DB.Exec("UPDATE accounts SET name = $1, currency = $2, balance = $3 WHERE id = $4 AND user_id = $5", req.Name, req.Currency, req.Balance, accID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}
	var acc models.Account
	h.DB.QueryRow("SELECT id, name, currency, balance FROM accounts WHERE id = $1 AND user_id = $2", accID, userID).Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Balance)
	c.JSON(http.StatusOK, acc)
}

// DeleteAccount deletes an account by id for the current user
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	accID := c.Param("id")
	_, err := h.DB.Exec("DELETE FROM accounts WHERE id = $1 AND user_id = $2", accID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}
	c.Status(http.StatusNoContent)
}
