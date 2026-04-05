package service

import (
	"context"

	"github.com/juanpblasi/go-template/internal/repository"
	"github.com/juanpblasi/go-template/pkg/errors"
	"github.com/juanpblasi/go-template/pkg/logger"
	"go.uber.org/zap"
)

type UserService interface {
	GetUser(ctx context.Context, id string) (*repository.User, error)
	CreateUser(ctx context.Context, name, email string) (*repository.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id string) (*repository.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("failed to get user", zap.String("id", id), zap.Error(err))
		return nil, errors.New(errors.ErrInternalError, "failed to get user")
	}
	if user == nil {
		return nil, errors.New(errors.ErrNotFound, "user not found")
	}

	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, name, email string) (*repository.User, error) {
	if name == "" || email == "" {
		return nil, errors.New(errors.ErrInvalidRequest, "name and email are required")
	}

	user := &repository.User{
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return nil, errors.New(errors.ErrInternalError, "failed to create user")
	}

	logger.Info("user created successfully", zap.String("id", user.ID))
	return user, nil
}
