package repositories

import (
	"database/sql"
	"log"

	"asset-dairy/models"
)

type AccountRepositoryInterface interface {
	ListAccounts(userID string) ([]models.Account, error)
	CreateAccount(userID string, acc *models.Account) error
	UpdateAccount(userID, accID string, req models.AccountUpdateRequest) (*models.Account, error)
	DeleteAccount(userID, accID string) error
}

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) ListAccounts(userID string) ([]models.Account, error) {
	rows, err := r.DB.Query("SELECT id, name, currency, balance FROM accounts WHERE user_id = $1", userID)
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

func (r *AccountRepository) CreateAccount(userID string, acc *models.Account) error {
	err := r.DB.QueryRow(
		"INSERT INTO accounts (id, user_id, name, currency, balance) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		acc.ID, userID, acc.Name, acc.Currency, acc.Balance,
	).Scan(&acc.ID)
	if err != nil {
		log.Println("Failed to create account:", err)
		return err
	}
	return nil
}

func (r *AccountRepository) UpdateAccount(userID, accID string, req models.AccountUpdateRequest) (*models.Account, error) {
	_, err := r.DB.Exec(
		"UPDATE accounts SET name = $1, currency = $2, balance = $3 WHERE id = $4 AND user_id = $5",
		req.Name, req.Currency, req.Balance, accID, userID,
	)
	if err != nil {
		log.Println("Failed to update account:", err)
		return nil, err
	}

	var acc models.Account
	err = r.DB.QueryRow(
		"SELECT id, name, currency, balance FROM accounts WHERE id = $1 AND user_id = $2",
		accID, userID,
	).Scan(&acc.ID, &acc.Name, &acc.Currency, &acc.Balance)
	if err != nil {
		log.Println("Failed to fetch updated account:", err)
		return nil, err
	}

	return &acc, nil
}

func (r *AccountRepository) DeleteAccount(userID, accID string) error {
	_, err := r.DB.Exec("DELETE FROM accounts WHERE id = $1 AND user_id = $2", accID, userID)
	if err != nil {
		log.Println("Failed to delete account:", err)
		return err
	}
	return nil
}
