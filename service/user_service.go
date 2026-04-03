package service

import (
	"context"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"go.uber.org/zap"
)

type (
	IUserService interface {
		GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error)
	}

	userService struct {
		userRepo   repository.IUserRepository
		jwtService jwt.IJWT
		logger     *zap.Logger
	}
)

func NewUserService(userRepo repository.IUserRepository, jwtService jwt.IJWT, logger *zap.Logger) *userService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (us *userService) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	data, found, err := us.userRepo.GetUserByID(ctx, nil, userID)
	if err != nil {
		us.logger.Error("failed to get user by id", zap.String("id", userID), zap.Error((err)))
		return nil, dto.ErrGetUserByID
	}

	if !found {
		us.logger.Warn("user not found", zap.String("id", userID))
		return nil, dto.ErrNotFound
	}

	return &dto.UserResponse{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
	}, nil
}
