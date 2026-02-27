package service

import (
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"go.uber.org/zap"
)

type (
	IUserService interface {
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
