package services

import (
	"errors"
	"nexus/backend/models"
	"nexus/backend/repositories"
)

var ErrUserHasServers = errors.New("cannot delete user with existing servers")

type UserService struct {
	userRepo   *repositories.UserRepository
	serverRepo *repositories.ServerRepository
}

func NewUserService(userRepo *repositories.UserRepository, serverRepo *repositories.ServerRepository) *UserService {
	return &UserService{userRepo: userRepo, serverRepo: serverRepo}
}

func (s *UserService) GetUserWithServers(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Delete(id uint) error {
	count := s.serverRepo.CountByUserID(id)
	if count > 0 {
		return ErrUserHasServers
	}
	return s.userRepo.Delete(id)
}

func (s *UserService) Update(user *models.User, password string) error {
	if password != "" {
		if err := user.HashPassword(password); err != nil {
			return err
		}
	}
	return s.userRepo.Update(user)
}
