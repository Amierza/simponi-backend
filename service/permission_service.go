package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"go.uber.org/zap"
)

type (
	IPermissionService interface {
		GetPermissions(ctx context.Context) ([]*dto.PermissionResponse, error)
		GetPermissionsWithPagination(ctx context.Context, req *response.PaginationRequest) (dto.PermissionPaginationResponse, error)
	}

	permissionService struct {
		permissionRepo repository.IPermissionRepository
		logger         *zap.Logger
		jwtService     jwt.IJWT
	}
)

func NewPermissionService(permissionRepo repository.IPermissionRepository, logger *zap.Logger, jwtService jwt.IJWT) *permissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		logger:         logger,
		jwtService:     jwtService,
	}
}

func mapToPermissionResponse(p *entity.Permission) *dto.PermissionResponse {
	return &dto.PermissionResponse{
		ID:       p.ID,
		Name:     p.Name,
		Endpoint: p.Endpoint,
		Method:   p.Method,
		Module:   p.Module,
	}
}

func (ps *permissionService) GetPermissions(ctx context.Context) ([]*dto.PermissionResponse, error) {
	datas, err := ps.permissionRepo.GetPermissions(ctx, nil)
	if err != nil {
		ps.logger.Error("failed to get permissions", zap.Error(err))
		return nil, fmt.Errorf("failed to get permissions: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to get permissions", zap.Int("count", len(datas)))

	var permissions []*dto.PermissionResponse
	for _, permission := range datas {
		permissions = append(permissions, mapToPermissionResponse(permission))
	}

	return permissions, nil
}
func (ps *permissionService) GetPermissionsWithPagination(ctx context.Context, req *response.PaginationRequest) (dto.PermissionPaginationResponse, error) {
	datas, err := ps.permissionRepo.GetPermissionsWithPagination(ctx, nil, req)
	if err != nil {
		ps.logger.Error("failed to get permissions", zap.Error(err))
		return dto.PermissionPaginationResponse{}, fmt.Errorf("failed to get permissions: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to get permissions", zap.Int64("count", datas.Count))

	var permissions []*dto.PermissionResponse
	for _, permission := range datas.Permissions {
		permissions = append(permissions, mapToPermissionResponse(permission))
	}

	return dto.PermissionPaginationResponse{
		Data:               permissions,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}
