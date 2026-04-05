package services

import (
	"errors"
	"nexus/backend/config"
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/utils"
)

var (
	ErrEmailTaken         = errors.New("email already in use")
	ErrUsernameTaken      = errors.New("username already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountSuspended   = errors.New("account is suspended")
)

type AuthService struct {
	userRepo *repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Register(req requests.RegisterRequest) (*models.User, string, error) {
	// Check email uniqueness
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, "", ErrEmailTaken
	}

	// Check username uniqueness
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, "", ErrUsernameTaken
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "client",
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, "", err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(*user, s.cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(req requests.LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, "", ErrInvalidCredentials
	}

	if user.Suspended {
		return nil, "", ErrAccountSuspended
	}

	token, err := utils.GenerateToken(*user, s.cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
