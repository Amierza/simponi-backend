package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IRoleRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IRoleRepository) error) error

		// CREATE
		CreateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error

		// READ
		GetRoleByID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) (*entity.Role, bool, error)
		GetRoleByName(ctx context.Context, tx *gorm.DB, name *string) (*entity.Role, bool, error)
		GetRoles(ctx context.Context, tx *gorm.DB) ([]*entity.Role, error)

		// UPDATE
		UpdateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error

		// DELETE
		DeleteRoleByID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) error
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

func (rr *roleRepository) RunInTransaction(ctx context.Context, fn func(txRepo IRoleRepository) error) error {
	return rr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &roleRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE
func (rr *roleRepository) CreateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if tx == nil {
		tx = rr.db
	}

	return tx.WithContext(ctx).Create(role).Error
}

// GET
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
func (rr *roleRepository) GetRoleByID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) (*entity.Role, bool, error) {
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
func (rr *roleRepository) GetRoles(ctx context.Context, tx *gorm.DB) ([]*entity.Role, error) {
	if tx == nil {
		tx = rr.db
	}

	var roles []*entity.Role
	if err := tx.WithContext(ctx).
		Model(&entity.Role{}).
		Preload("Users").
		Preload("RolePermissions.Permission").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// UPDATE
func (rr *roleRepository) UpdateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if tx == nil {
		tx = rr.db
	}

	return tx.WithContext(ctx).Model(&entity.Role{}).Where("id = ?", role.ID).Updates(&role).Error
}

// DELETE
func (rr *roleRepository) DeleteRoleByID(ctx context.Context, tx *gorm.DB, roleID *uuid.UUID) error {
	if tx == nil {
		tx = rr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &roleID).Delete(&entity.Role{}).Error
}
