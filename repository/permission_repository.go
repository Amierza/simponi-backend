package repository

import (
	"context"

	"gorm.io/gorm"
)

type (
	IPermissionRepository interface {
		CheckPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error)
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

func (pr *permissionRepository) CheckPermissionByName(ctx context.Context, roleID, permissionName string) (bool, error) {
	var count int64
	err := pr.db.WithContext(ctx).
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
