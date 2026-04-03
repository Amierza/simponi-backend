package service

import (
	"context"

	"github.com/Amierza/simponi-backend/repository"
)

type (
	IPermissionService interface {
		HasPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error)
	}

	permissionService struct {
		permissionRepo repository.IPermissionRepository
	}
)

func NewPermissionService(permissionRepo repository.IPermissionRepository) *permissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

func (ps *permissionService) HasPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error) {
	return ps.permissionRepo.CheckPermissionByName(ctx, roleID, permissionName)
}
