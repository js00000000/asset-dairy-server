package services

import (
	"asset-dairy/models"
	"asset-dairy/repositories"
)

type TradeServiceInterface interface {
	ListTrades(userID string) ([]models.Trade, error)
}

type TradeService struct {
	repo repositories.TradeRepositoryInterface
}

// NewTradeService creates a new TradeService instance with a repository
func NewTradeService(repo repositories.TradeRepositoryInterface) *TradeService {
	return &TradeService{repo: repo}
}

// ListTrades retrieves all trades for a given user
func (s *TradeService) ListTrades(userID string) ([]models.Trade, error) {
	return s.repo.ListTrades(userID)
}
