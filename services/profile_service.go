package services

import (
	"asset-dairy/models"
	"asset-dairy/repositories"
)

type ProfileServiceInterface interface {
	GetProfile(userID string) (*models.ProfileResponse, error)
	ChangePassword(userID string, currentPassword, newPassword string) error
	UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error)
}

type ProfileService struct {
	repo repositories.ProfileRepositoryInterface
}

// NewProfileService creates a new ProfileService instance with a repository
func NewProfileService(repo repositories.ProfileRepositoryInterface) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}

func (s *ProfileService) GetProfile(userID string) (*models.ProfileResponse, error) {
	return s.repo.GetProfile(userID)
}

func (s *ProfileService) ChangePassword(userID string, currentPassword, newPassword string) error {
	return s.repo.ChangePassword(userID, currentPassword, newPassword)
}

func (s *ProfileService) UpdateProfile(userID string, req *models.UserUpdateRequest) (*models.User, error) {
	return s.repo.UpdateProfile(userID, req)
}


