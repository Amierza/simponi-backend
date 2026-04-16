package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IRoleService interface {
		CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
		GetRoles(ctx context.Context, req *response.PaginationRequest) (dto.RolePaginationResponse, error)
		GetRoleByID(ctx context.Context, roleID *uuid.UUID) (*dto.RoleResponse, error)
		UpdateRole(ctx context.Context, roleID *uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
		DeleteRoleByID(ctx context.Context, roleID *uuid.UUID) error
	}

	roleService struct {
		roleRepo           repository.IRoleRepository
		permissionRepo     repository.IPermissionRepository
		rolePermissionRepo repository.IRolePermissionRepository
		logger             *zap.Logger
		jwtService         jwt.IJWT
	}
)

func NewRoleService(roleRepo repository.IRoleRepository, permissionRepo repository.IPermissionRepository, rolePermissionRepo repository.IRolePermissionRepository, logger *zap.Logger, jwtService jwt.IJWT) *roleService {
	return &roleService{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		logger:             logger,
		jwtService:         jwtService,
	}
}

func mapToRoleResponse(r *entity.Role, p []*entity.Permission) *dto.RoleResponse {
	res := &dto.RoleResponse{
		ID:   r.ID,
		Name: r.Name,
	}
	for _, permission := range p {
		res.Permissions = append(res.Permissions, dto.PermissionResponse{
			ID:       permission.ID,
			Name:     permission.Name,
			Endpoint: permission.Endpoint,
			Method:   permission.Method,
			Module:   permission.Module,
		})
	}
	return res
}

func (rs *roleService) CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	_, found, err := rs.roleRepo.GetRoleByName(ctx, nil, &req.Name)
	if err != nil {
		rs.logger.Error("failed to get role by name", zap.String("name", req.Name), zap.Error(err))
		return nil, fmt.Errorf("failed to get role by name: %w", dto.ErrInternal)
	}
	if found {
		rs.logger.Warn("role already exists", zap.String("name", req.Name))
		return nil, fmt.Errorf("role already exists: %w", dto.ErrAlreadyExists)
	}

	var permissions []*entity.Permission
	for _, permissionID := range req.PermissionsIDs {
		permission, found, err := rs.permissionRepo.GetPermissionByID(ctx, nil, permissionID)
		if err != nil {
			rs.logger.Error("failed to get permission by id", zap.String("id", permission.ID.String()), zap.Error(err))
			return nil, fmt.Errorf("failed to get permission by id: %w", dto.ErrInternal)
		}
		if !found {
			rs.logger.Warn("permission not found", zap.String("name", permission.Name))
			return nil, fmt.Errorf("permission not found: %w", dto.ErrNotFound)
		}
		permissions = append(permissions, permission)
	}

	newID := uuid.New()
	newRole := &entity.Role{
		ID:   newID,
		Name: req.Name,
	}
	err = rs.roleRepo.CreateRole(ctx, nil, newRole)
	if err != nil {
		rs.logger.Error("failed to create role", zap.Error(err))
		return nil, fmt.Errorf("failed to create role: %w", dto.ErrInternal)
	}
	rs.logger.Info("success to create role", zap.String("id", newRole.ID.String()))

	for _, permission := range permissions {
		_, err := rs.rolePermissionRepo.CreateRolePermission(ctx, nil, &entity.RolePermission{
			ID:           uuid.New(),
			RoleID:       &newRole.ID,
			PermissionID: &permission.ID,
		})
		if err != nil {
			rs.logger.Error("failed to create role permission", zap.Error(err))
			return nil, fmt.Errorf("failed to create role permission: %w", dto.ErrInternal)
		}
	}
	rs.logger.Info("success to create role permissions")

	return mapToRoleResponse(newRole, permissions), nil
}

func (rs *roleService) GetRoles(ctx context.Context, req *response.PaginationRequest) (dto.RolePaginationResponse, error) {
	datas, err := rs.roleRepo.GetRoles(ctx, nil, req)
	if err != nil {
		rs.logger.Error("failed to get roles", zap.Error(err))
		return dto.RolePaginationResponse{}, fmt.Errorf("failed to get roles: %w", dto.ErrInternal)
	}
	rs.logger.Info("success to get roles", zap.Int64("count", datas.Count))

	var roles []*dto.RoleResponse
	for _, role := range datas.Roles {
		rawPermissions, err := rs.permissionRepo.GetPermissionsByRoleID(ctx, nil, &role.ID)
		if err != nil {
			rs.logger.Error("failed to get permissions by role id", zap.Error(err))
			return dto.RolePaginationResponse{}, fmt.Errorf("failed to get permissions by role id: %w", dto.ErrInternal)
		}

		roles = append(roles, mapToRoleResponse(role, rawPermissions))
	}

	return dto.RolePaginationResponse{
		Data:               roles,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (rs *roleService) GetRoleByID(ctx context.Context, roleID *uuid.UUID) (*dto.RoleResponse, error) {
	role, found, err := rs.roleRepo.GetRoleByID(ctx, nil, roleID)
	if err != nil {
		rs.logger.Error("failed to get role by ID", zap.String("roleID", roleID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get role ID: %w", dto.ErrInternal)
	}
	if !found {
		rs.logger.Warn("role not found", zap.String("roleID", roleID.String()))
		return nil, fmt.Errorf("role not found: %v", dto.ErrNotFound)
	}
	rs.logger.Info("success to get role by id", zap.String("id", roleID.String()))

	rawPermissions, err := rs.permissionRepo.GetPermissionsByRoleID(ctx, nil, &role.ID)
	if err != nil {
		rs.logger.Error("failed to get permissions by role id", zap.Error(err))
		return nil, fmt.Errorf("failed to get permissions by role id: %w", dto.ErrInternal)
	}
	rs.logger.Info("success to get permission by role id", zap.String("role id", roleID.String()))

	return mapToRoleResponse(role, rawPermissions), nil
}

func (rs *roleService) UpdateRole(ctx context.Context, roleID *uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, found, err := rs.roleRepo.GetRoleByID(ctx, nil, roleID)
	if err != nil {
		rs.logger.Error("failed to get role by ID", zap.String("roleID", roleID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get role ID: %w", dto.ErrInternal)
	}
	if !found {
		rs.logger.Warn("role not found", zap.String("roleID", roleID.String()))
		return nil, fmt.Errorf("role not found: %v", dto.ErrNotFound)
	}
	role.Name = req.Name

	var permissions []*entity.Permission
	for _, permissionID := range req.PermissionIDs {
		permission, found, err := rs.permissionRepo.GetPermissionByID(ctx, nil, permissionID)
		if err != nil {
			rs.logger.Error("failed to get permission by name", zap.String("name", permission.Name), zap.Error(err))
			return nil, fmt.Errorf("failed to get permission by name: %w", dto.ErrInternal)
		}
		if !found {
			rs.logger.Warn("permission not found", zap.String("name", permission.Name))
			return nil, fmt.Errorf("permission not found: %w", dto.ErrNotFound)
		}
		permissions = append(permissions, permission)
	}

	err = rs.roleRepo.RunInTransaction(ctx, func(txRepo repository.IRoleRepository) error {
		err = txRepo.UpdateRole(ctx, nil, role)
		if err != nil {
			rs.logger.Error("failed to update role", zap.String("id", roleID.String()), zap.Error(err))
			return fmt.Errorf("failed to update role: %w", dto.ErrInternal)
		}

		if err := rs.rolePermissionRepo.DeleteRolePermissionsByRoleID(ctx, nil, &role.ID); err != nil {
			return fmt.Errorf("failed to delete role permissions by role id: %w", dto.ErrInternal)
		}

		for _, permission := range permissions {
			_, err := rs.rolePermissionRepo.CreateRolePermission(ctx, nil, &entity.RolePermission{
				ID:           uuid.New(),
				RoleID:       &role.ID,
				PermissionID: &permission.ID,
			})
			if err != nil {
				rs.logger.Error("failed to create role permission", zap.Error(err))
				return fmt.Errorf("failed to create role permission: %w", dto.ErrInternal)
			}
		}

		return nil
	})
	if err != nil {
		return &dto.RoleResponse{}, err
	}

	return mapToRoleResponse(role, permissions), nil
}

func (rs *roleService) DeleteRoleByID(ctx context.Context, roleID *uuid.UUID) error {
	role, found, err := rs.roleRepo.GetRoleByID(ctx, nil, roleID)
	if err != nil {
		rs.logger.Error("failed to get role by ID", zap.String("roleID", roleID.String()), zap.Error(err))
		return fmt.Errorf("failed to get role ID: %w", dto.ErrInternal)
	}
	if !found {
		rs.logger.Warn("role not found", zap.String("roleID", roleID.String()))
		return fmt.Errorf("role not found: %v", dto.ErrNotFound)
	}

	if err := rs.rolePermissionRepo.DeleteRolePermissionsByRoleID(ctx, nil, &role.ID); err != nil {
		return fmt.Errorf("failed to delete role permissions by role id: %w", dto.ErrInternal)
	}

	if err := rs.roleRepo.DeleteRoleByID(ctx, nil, roleID); err != nil {
		rs.logger.Error("failed to delete role by id", zap.String("roleID", roleID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete role by id: %w", dto.ErrInternal)
	}

	return nil
}
