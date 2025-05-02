package services

import (
	"database/sql"
	"log"

	"asset-dairy/models"

	"github.com/google/uuid"
)

type AccountServiceInterface interface {
	ListAccounts(userID string) ([]models.Account, error)
	CreateAccount(userID string, req models.AccountCreateRequest) (*models.Account, error)
	UpdateAccount(userID, accID string, req models.AccountUpdateRequest) (*models.Account, error)
	DeleteAccount(userID, accID string) error
}

type AccountService struct {
	DB *sql.DB
}

func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{DB: db}
}

func (s *AccountService) ListAccounts(userID string) ([]models.Account, error) {
	rows, err := s.DB.Query("SELECT id, name, currency, balance FROM accounts WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Failed to fetch accounts:", err)
		return nil, err
	}
	defer rows.Close()

	accounts := []models.Account{}
	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *AccountService) CreateAccount(userID string, req models.AccountCreateRequest) (*models.Account, error) {
	id := uuid.New().String()
	acc := &models.Account{
		ID:       id,
		Name:     req.Name,
		Currency: req.Currency,
		Balance:  req.Balance,
	}

	err := s.DB.QueryRow(
		"INSERT INTO accounts (id, user_id, name, currency, balance) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		id, userID, req.Name, req.Currency, req.Balance,
	).Scan(&acc.ID)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) UpdateAccount(userID, accID string, req models.AccountUpdateRequest) (*models.Account, error) {
	_, err := s.DB.Exec(
		"UPDATE accounts SET name = $1, currency = $2, balance = $3 WHERE id = $4 AND user_id = $5",
		req.Name, req.Currency, req.Balance, accID, userID,
	)
	if err != nil {
		return nil, err
	}

	var acc models.Account
	err = s.DB.QueryRow(
		"SELECT id, name, currency, balance FROM accounts WHERE id = $1 AND user_id = $2",
		accID, userID,
	).Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Balance)

	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (s *AccountService) DeleteAccount(userID, accID string) error {
	_, err := s.DB.Exec("DELETE FROM accounts WHERE id = $1 AND user_id = $2", accID, userID)
	return err
}
