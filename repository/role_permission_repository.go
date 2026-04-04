package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IRolePermissionRepository interface {
		// CREATE
		CreateRolePermission(ctx context.Context, tx *gorm.DB, rolePermission *entity.RolePermission) (*entity.RolePermission, error)

		// READ
		CheckRolePermissionByPermissionName(ctx context.Context, tx *gorm.DB, roleID, permissionName string) (bool, error)

		// DELETE
		DeleteRolePermissionsByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) error
	}

	rolePermissionRepository struct {
		db *gorm.DB
	}
)

func NewRolePermissionRepository(db *gorm.DB) *rolePermissionRepository {
	return &rolePermissionRepository{
		db: db,
	}
}

// CREATE
func (rpr *rolePermissionRepository) CreateRolePermission(ctx context.Context, tx *gorm.DB, rolePermission *entity.RolePermission) (*entity.RolePermission, error) {
	if tx == nil {
		tx = rpr.db
	}

	if err := tx.WithContext(ctx).Create(rolePermission).Error; err != nil {
		return nil, err
	}

	return rolePermission, nil
}

// GET
func (rpr *rolePermissionRepository) CheckRolePermissionByPermissionName(ctx context.Context, tx *gorm.DB, roleID, permissionName string) (bool, error) {
	if tx == nil {
		tx = rpr.db
	}

	var count int64
	err := rpr.db.WithContext(ctx).
		Table("role_permissions rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Where("p.name = ?", permissionName).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DELETE
func (rpr *rolePermissionRepository) DeleteRolePermissionsByRoleID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) error {
	if tx == nil {
		tx = rpr.db
	}

	return tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&entity.RolePermission{}).Error
}
