package service

import (
	"context"

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

	authSevice struct {
		authRepo repository.IAuthRepository
		logger   *zap.Logger
		jwt      jwt.IJWT
	}
)

func NewAuthService(authRepo repository.IAuthRepository, logger *zap.Logger, jwt jwt.IJWT) *authSevice {
	return &authSevice{
		authRepo: authRepo,
		logger:   logger,
		jwt:      jwt,
	}
}

func (as *authSevice) SignIn(ctx context.Context, req dto.SignInRequest) (dto.SignInResponse, error) {
	user, found, err := as.authRepo.GetUserByEmail(ctx, nil, &req.Email)
	if err != nil {
		as.logger.Error("failed to get user by email", zap.String("email", req.Email), zap.Error(err))
		return dto.SignInResponse{}, dto.ErrGetUserByEmail
	}
	if !found {
		as.logger.Warn("user not found", zap.String("email", req.Email))
		return dto.SignInResponse{}, dto.ErrNotFound
	}

	checkPassword, err := helper.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		as.logger.Error("incorrect credentials", zap.String("passwrd", req.Password), zap.Error(err))
		return dto.SignInResponse{}, dto.ErrIncorrectPassword
	}

	accessToken, refreshToken, err := as.jwt.GenerateToken(user.ID.String())
	if err != nil {
		as.logger.Error("failed to generate access token and refresh token", zap.String("email", req.Email), zap.Error(err))
		return dto.SignInResponse{}, dto.ErrGenerateAccessAndRefreshToken
	}

	as.logger.Info("Sign In Success", zap.String("email", req.Email))

	return dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *authSevice) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error) {
	_, err := as.jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		as.logger.Error("invalid token", zap.String("refresh_token", req.RefreshToken), zap.Error(err))
		return dto.RefreshTokenResponse{}, dto.ErrValidateToken
	}

	userID, err := as.jwt.GetUserIDByToken(req.RefreshToken)
	if err != nil {
		as.logger.Error("failed to get user id by token", zap.String("refresh_token", req.RefreshToken), zap.Error(err))
		return dto.RefreshTokenResponse{}, dto.ErrGetUserIDFromToken
	}

	accessToken, _, err := as.jwt.GenerateToken(userID)
	if err != nil {
		as.logger.Error("failed to generate token", zap.String("user_id", userID), zap.Error(err))
		return dto.RefreshTokenResponse{}, dto.ErrGenerateAccessToken
	}

	as.logger.Info("refresh token sucess", zap.String("access_token", accessToken))

	return dto.RefreshTokenResponse{AccessToken: accessToken}, nil
}
