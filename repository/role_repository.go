package repository

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IRoleRepository interface {
		// CREATE
		CreateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error

		// READ
		GetRoleByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) (*entity.Role, bool, error)
		GetRoleByName(ctx context.Context, tx *gorm.DB, name *string) (*entity.Role, bool, error)
		GetRoles(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.RolePaginationRepositoryResponse, error)

		// UPDATE
		UpdateRoleByRoleID(ctx context.Context, tx *gorm.DB, role *entity.Role) error

		// DELETE
		DeleteRoleByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) error
	}

	roleRepository struct {
		db *gorm.DB
	}
)

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

// CREATE
func (rr *roleRepository) CreateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if tx == nil {
		tx = rr.db
	}

	return tx.WithContext(ctx).Create(role).Error
}

// GET
func (rr *roleRepository) GetRoleByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) (*entity.Role, bool, error) {
	if tx == nil {
		tx = rr.db
	}

	var role *entity.Role
	err := tx.WithContext(ctx).
		Model(&entity.Role{}).
		Preload("Users").
		Preload("RolePermissions").
		Where("id = ?", roleID).
		First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return role, true, nil
}
func (rr *roleRepository) GetRoleByName(ctx context.Context, tx *gorm.DB, name *string) (*entity.Role, bool, error) {
	if tx == nil {
		tx = rr.db
	}

	var role *entity.Role
	err := tx.WithContext(ctx).
		Model(&entity.Role{}).
		Preload("Users").
		Preload("RolePermissions").
		Where("name = ?", name).
		First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return role, true, nil
}
func (rr *roleRepository) GetRoles(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.RolePaginationRepositoryResponse, error) {
	if tx == nil {
		tx = rr.db
	}

	var roles []*entity.Role
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Role{}).
		Preload("Users").
		Preload("RolePermissions.Permission")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(`
			LOWER(roles.name) LIKE ?
			OR roles.id IN (
				SELECT rp.role_id
				FROM role_permissions rp
				JOIN permissions p ON p.id = rp.permission_id
				WHERE LOWER(p.module) LIKE ?
			)
		`, searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.RolePaginationRepositoryResponse{}, err
	}

	if err := query.Order("roles.created_at DESC").Scopes(response.Paginate(req.Page, req.PerPage)).Find(&roles).Error; err != nil {
		return dto.RolePaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.RolePaginationRepositoryResponse{
		Roles: roles,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

// UPDATE
func (rr *roleRepository) UpdateRoleByRoleID(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if tx == nil {
		tx = rr.db
	}

	return tx.WithContext(ctx).Model(&entity.Role{}).Where("id = ?", role.ID).Updates(&role).Error
}

// DELETE
func (rr *roleRepository) DeleteRoleByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) error {
	if tx == nil {
		tx = rr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &roleID).Delete(&entity.Role{}).Error
}
