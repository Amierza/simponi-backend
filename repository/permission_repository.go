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
	IPermissionRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IPermissionRepository) error) error

		// READ
		GetPermissionByPermissionID(ctx context.Context, tx *gorm.DB, permissionID *uuid.UUID) (*entity.Permission, bool, error)
		GetPermissions(ctx context.Context, tx *gorm.DB) ([]*entity.Permission, error)
		GetPermissionsWithPagination(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.PermissionPaginationRepositoryResponse, error)
		GetPermissionsByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) ([]*entity.Permission, error)
	}

	permissionRepository struct {
		db *gorm.DB
	}
)

func NewPermissionRepository(db *gorm.DB) *permissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (pr *permissionRepository) RunInTransaction(ctx context.Context, fn func(txRepo IPermissionRepository) error) error {
	return pr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &permissionRepository{db: tx}
		return fn(txRepo)
	})
}

// READ
func (pr *permissionRepository) GetPermissionByPermissionID(ctx context.Context, tx *gorm.DB, permissionID *uuid.UUID) (*entity.Permission, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	permission := new(entity.Permission)
	err := tx.WithContext(ctx).
		Preload("RolePermissions").
		Where("id = ?", permissionID).
		Take(permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return permission, true, nil
}
func (pr *permissionRepository) GetPermissionsByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) ([]*entity.Permission, error) {
	if tx == nil {
		tx = pr.db
	}

	var permissions []*entity.Permission
	err := tx.WithContext(ctx).
		Table("permissions p").
		Select("p.*").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Where("rp.role_id = ? AND rp.deleted_at IS NULL", *roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0 {
		return nil, nil
	}

	return permissions, nil
}
func (pr *permissionRepository) GetPermissions(ctx context.Context, tx *gorm.DB) ([]*entity.Permission, error) {
	if tx == nil {
		tx = pr.db
	}

	var permissions []*entity.Permission

	query := tx.WithContext(ctx).
		Model(&entity.Permission{}).
		Preload("RolePermissions")

	if err := query.Order(`"created_at" DESC`).Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}
func (pr *permissionRepository) GetPermissionsWithPagination(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.PermissionPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var permissions []*entity.Permission
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Permission{})

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(endpoint) LIKE ? OR LOWER(method) LIKE ? OR LOWER(module) LIKE ?", searchValue, searchValue, searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.PermissionPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&permissions).Error; err != nil {
		return dto.PermissionPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.PermissionPaginationRepositoryResponse{
		Permissions: permissions,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
