package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IImpersonateService interface {
		StartImpersonate(ctx context.Context, adminID, targetUserID string) (*dto.ImpersonateResponse, error)
		StopImpersonate(ctx context.Context, claims *jwt.CustomClaims) (*dto.ImpersonateResponse, error)
	}

	impersonateService struct {
		userRepo       repository.IUserRepository
		permissionRepo repository.IPermissionRepository
		logger         *zap.Logger
		jwt            jwt.IJWT
	}
)

func NewImpersonateService(userRepo repository.IUserRepository, permissionRepo repository.IPermissionRepository, logger *zap.Logger, jwt jwt.IJWT) *impersonateService {
	return &impersonateService{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
		logger:         logger,
		jwt:            jwt,
	}
}

func mapPermissions(perms []*entity.Permission) []string {
	var result []string
	for _, p := range perms {
		result = append(result, p.Name)
	}
	return result
}

func (is *impersonateService) StartImpersonate(ctx context.Context, adminID, targetUserID string) (*dto.ImpersonateResponse, error) {
	// parse UUID
	adminUUID, err := uuid.Parse(adminID)
	if err != nil {
		return nil, fmt.Errorf("invalid admin id: %w", dto.ErrBadRequest)
	}

	// ambil admin
	admin, found, err := is.userRepo.GetUserByID(ctx, nil, &adminUUID)
	if err != nil {
		is.logger.Error("failed to get admin", zap.Error(err))
		return nil, fmt.Errorf("failed to get admin: %w", dto.ErrInternal)
	}
	if !found {
		return nil, fmt.Errorf("admin not found: %w", dto.ErrNotFound)
	}

	// validasi role
	if admin.Role.Name != "superadmin" {
		is.logger.Warn("unauthorized impersonate attempt", zap.String("admin_id", adminID))
		return nil, fmt.Errorf("access denied: %w", dto.ErrForbidden)
	}

	// ambil target user
	targetUUID, err := uuid.Parse(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid target user id: %w", dto.ErrBadRequest)
	}

	targetUser, found, err := is.userRepo.GetUserByID(ctx, nil, &targetUUID)
	if err != nil {
		is.logger.Error("failed to get target user", zap.Error(err))
		return nil, fmt.Errorf("failed to get target user: %w", dto.ErrInternal)
	}
	if !found {
		return nil, fmt.Errorf("target user not found: %w", dto.ErrNotFound)
	}

	// ambil permissions target user
	perms, err := is.permissionRepo.GetPermissionsByRoleID(ctx, nil, targetUser.RoleID)
	if err != nil {
		is.logger.Error("failed to get permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to get permissions: %w", dto.ErrInternal)
	}

	permissionNames := mapPermissions(perms)

	// generate token
	duration := 15 * time.Minute

	token, err := is.jwt.GenerateImpersonateToken(
		targetUser.ID.String(),
		targetUser.RoleID.String(),
		admin.ID.String(),
		permissionNames,
		duration,
	)
	if err != nil {
		is.logger.Error("failed to generate impersonate token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", dto.ErrInternal)
	}

	is.logger.Info("impersonation started",
		zap.String("admin_id", adminID),
		zap.String("target_user_id", targetUserID),
	)

	return &dto.ImpersonateResponse{
		AccessToken: token,
	}, nil
}

func (is *impersonateService) StopImpersonate(ctx context.Context, claims *jwt.CustomClaims) (*dto.ImpersonateResponse, error) {
	// validasi impersonate
	if !claims.IsImpersonating {
		return nil, fmt.Errorf("not impersonating: %w", dto.ErrBadRequest)
	}

	if claims.OriginalUserID == "" {
		return nil, fmt.Errorf("invalid original user: %w", dto.ErrBadRequest)
	}

	originalUUID, err := uuid.Parse(claims.OriginalUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid original user id: %w", dto.ErrBadRequest)
	}

	// ambil user asli
	user, found, err := is.userRepo.GetUserByID(ctx, nil, &originalUUID)
	if err != nil {
		is.logger.Error("failed to get original user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", dto.ErrInternal)
	}
	if !found {
		return nil, fmt.Errorf("user not found: %w", dto.ErrNotFound)
	}

	// ambil permissions user asli
	perms, err := is.permissionRepo.GetPermissionsByRoleID(ctx, nil, user.RoleID)
	if err != nil {
		is.logger.Error("failed to get permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to get permissions: %w", dto.ErrInternal)
	}

	permissionNames := mapPermissions(perms)

	// generate token normal
	duration := 5 * time.Minute

	token, err := is.jwt.GenerateToken(
		user.ID.String(),
		user.RoleID.String(),
		permissionNames,
		duration,
	)
	if err != nil {
		is.logger.Error("failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", dto.ErrInternal)
	}

	is.logger.Info("impersonation stopped",
		zap.String("original_user_id", claims.OriginalUserID),
	)

	return &dto.ImpersonateResponse{
		AccessToken: token,
	}, nil
}
