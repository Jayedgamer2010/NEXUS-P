package services

import (
	"errors"
	"fmt"

	"nexus/backend/config"
	"nexus/backend/models"
	"nexus/backend/repositories"
	"nexus/backend/requests"
	"nexus/backend/utils"
)

var (
	ErrEmailTaken    = errors.New("email already in use")
	ErrUsernameTaken = errors.New("username already in use")
)

type UserService struct {
	repo  *repositories.UserRepository
	cfg   *config.Config
}

func NewUserService(repo *repositories.UserRepository, cfg *config.Config) *UserService {
	return &UserService{repo: repo, cfg: cfg}
}

func (s *UserService) Register(req requests.RegisterRequest) (*models.User, string, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, "", ErrEmailTaken
	}

	existing, err = s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check username: %w", err)
	}
	if existing != nil {
		return nil, "", ErrUsernameTaken
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "client",
		Coins:    0,
	}
	if err := user.HashPassword(req.Password); err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := utils.GenerateToken(user.ID, user.UUID, user.Role, s.cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *UserService) Login(req requests.LoginRequest) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID, user.UUID, user.Role, s.cfg)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

func (s *UserService) All(page, perPage int) ([]models.User, int64, error) {
	return s.repo.All(page, perPage)
}

func (s *UserService) FindByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Create(req requests.CreateUserRequest) (*models.User, error) {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, ErrEmailTaken
	}

	existing, _ = s.repo.FindByUsername(req.Username)
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      req.Role,
		Coins:     req.Coins,
		RootAdmin: req.Role == "admin",
	}
	if err := user.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(user *models.User, req requests.UpdateUserRequest) error {
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil && *req.Password != "" {
		if err := user.HashPassword(*req.Password); err != nil {
			return err
		}
	}
	if req.Role != nil {
		user.Role = *req.Role
		user.RootAdmin = *req.Role == "admin"
	}
	if req.Coins != nil {
		user.Coins = *req.Coins
	}
	if req.RootAdmin != nil {
		user.RootAdmin = *req.RootAdmin
	}

	return s.repo.Update(user)
}

func (s *UserService) Delete(user *models.User) error {
	return s.repo.Delete(user)
}
