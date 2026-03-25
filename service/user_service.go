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
		GetProfile(ctx context.Context) (*dto.UserResponse, error)
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

func (us *userService) GetProfile(ctx context.Context) (*dto.UserResponse, error){
	token := ctx.Value("Authorization").(string)
	userIDString, err := us.jwtService.GetUserIDByToken(token)
	if err != nil {
		us.logger.Error("failed to get user_id by token", zap.String("id", userIDString), zap.Error(err))
		return &dto.UserResponse{}, dto.ErrGetUserIDFromToken
	}

	data, found, err := us.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		us.logger.Error("failed to get user by id", zap.String("id", userIDString), zap.Error((err)))
		return nil, dto.ErrGetUserByID
	}

	if !found{
		us.logger.Warn("user not found", zap.String("id", userIDString))
		return nil, dto.ErrNotFound
	}

	return &dto.UserResponse{
		ID: data.ID,
		Name: data.Name,
		Email: data.Email,
	}, nil
}
