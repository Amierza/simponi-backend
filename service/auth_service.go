package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/helper"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"go.uber.org/zap"
)

type (
	IAuthService interface {
		SignIn(ctx context.Context, req dto.SignInRequest) (dto.SignInResponse, error)
		RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error)
	}

	authService struct {
		userRepo       repository.IUserRepository
		permissionRepo repository.IPermissionRepository
		logger         *zap.Logger
		jwt            jwt.IJWT
	}
)

func NewAuthService(userRepo repository.IUserRepository, permissionRepo repository.IPermissionRepository, logger *zap.Logger, jwt jwt.IJWT) *authService {
	return &authService{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
		logger:         logger,
		jwt:            jwt,
	}
}

func (as *authService) SignIn(ctx context.Context, req dto.SignInRequest) (dto.SignInResponse, error) {
	user, found, err := as.userRepo.GetUserByEmail(ctx, nil, &req.Email)
	if err != nil {
		as.logger.Error("failed to get user by email", zap.String("email", req.Email), zap.Error(err))
		return dto.SignInResponse{}, fmt.Errorf("failed to get user by email: %w", dto.ErrInternal)
	}
	if !found {
		as.logger.Warn("user not found", zap.String("email", req.Email))
		return dto.SignInResponse{}, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	checkPassword, err := helper.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		as.logger.Error("incorrect password", zap.String("password", req.Password), zap.Error(err))
		return dto.SignInResponse{}, fmt.Errorf("incorrect password: %w", dto.ErrBadRequest)
	}

	permissions, err := as.permissionRepo.GetPermissionsByRoleID(ctx, nil, user.RoleID)
	if err != nil {
		as.logger.Error("failed to get permissions by role id", zap.String("role id", user.RoleID.String()), zap.Error(err))
		return dto.SignInResponse{}, fmt.Errorf("failed to get permissions by role id: %w", dto.ErrInternal)
	}

	accessToken, err := as.jwt.GenerateToken(user.ID.String(), user.RoleID.String(), mapPermissions(permissions), 5*time.Minute)
	if err != nil {
		as.logger.Error("failed to generate access token", zap.String("email", req.Email), zap.Error(err))
		return dto.SignInResponse{}, fmt.Errorf("failed to generate access token: %w", dto.ErrInternal)
	}

	refreshToken, err := as.jwt.GenerateToken(user.ID.String(), user.RoleID.String(), mapPermissions(permissions), 7*24*time.Hour)
	if err != nil {
		as.logger.Error("failed to generate refresh token", zap.String("email", req.Email), zap.Error(err))
		return dto.SignInResponse{}, fmt.Errorf("failed to generate refresh token: %w", dto.ErrInternal)
	}

	as.logger.Info("Sign In Success", zap.String("email", req.Email))

	return dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error) {
	claims, err := as.jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		as.logger.Error("invalid access token", zap.String("refresh_token", req.RefreshToken), zap.Error(err))
		return dto.RefreshTokenResponse{}, fmt.Errorf("invalid access token: %w", dto.ErrInternal)
	}

	userID := claims.UserID
	roleID := claims.RoleID
	permissions := claims.Permissions

	accessToken, err := as.jwt.GenerateToken(userID, roleID, permissions, 5*time.Minute)
	if err != nil {
		as.logger.Error("failed to generate token", zap.String("user_id", userID), zap.String("role_id", roleID), zap.Error(err))
		return dto.RefreshTokenResponse{}, fmt.Errorf("failed to generate token: %w", dto.ErrInternal)
	}

	as.logger.Info("refresh token sucess", zap.String("access_token", accessToken))

	return dto.RefreshTokenResponse{AccessToken: accessToken}, nil
}
