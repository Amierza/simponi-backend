package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IStoreCredentialRepository interface {
		// CREATE
		CreateStoreCredential(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) error

		// READ

		// UPDATE
		UpdateStoreCredential(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) error

		// DELETE
		DeleteStoreCredentialByID(ctx context.Context, tx *gorm.DB, id *uuid.UUID) error
	}

	storeCredentialRepository struct {
		db *gorm.DB
	}
)

func NewStoreCredentialRepository(db *gorm.DB) *storeCredentialRepository {
	return &storeCredentialRepository{
		db: db,
	}
}

// CREATE
func (scr *storeCredentialRepository) CreateStoreCredential(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) error {
	if tx == nil {
		tx = scr.db
	}

	return tx.WithContext(ctx).Create(storeCredential).Error
}

// READ

// UPDATE
func (scr *storeCredentialRepository) UpdateStoreCredential(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) error {
	if tx == nil {
		tx = scr.db
	}

	return tx.WithContext(ctx).Model(&entity.StoreCredential{}).Where("id = ?", storeCredential.ID).Updates(storeCredential).Error
}

// DELETE
func (r *storeCredentialRepository) DeleteStoreCredentialByID(ctx context.Context, tx *gorm.DB, id *uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.StoreCredential{}).Error
}
