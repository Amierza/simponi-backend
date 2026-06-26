package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IStoreCredentialRepository interface {
		// CREATE
		CreateStoreCredential(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) error

		// READ
		UpsertStoreCredentialByStorePlatformID(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) (*entity.StoreCredential, error)

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
func (scr *storeCredentialRepository) UpsertStoreCredentialByStorePlatformID(ctx context.Context, tx *gorm.DB, storeCredential *entity.StoreCredential) (*entity.StoreCredential, error) {
	if tx == nil {
		tx = scr.db
	}

	var existing entity.StoreCredential
	err := tx.WithContext(ctx).
		Unscoped().
		Where("store_platform_id = ?", storeCredential.StorePlatformID).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.WithContext(ctx).Create(storeCredential).Error; err != nil {
			return nil, err
		}

		return storeCredential, nil
	}
	if err != nil {
		return nil, err
	}

	err = tx.WithContext(ctx).
		Unscoped().
		Model(&entity.StoreCredential{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"access_token":  storeCredential.AccessToken,
			"refresh_token": storeCredential.RefreshToken,
			"expires_at":    storeCredential.ExpiresAt,
			"deleted_at":    nil,
		}).Error
	if err != nil {
		return nil, err
	}

	existing.AccessToken = storeCredential.AccessToken
	existing.RefreshToken = storeCredential.RefreshToken
	existing.ExpiresAt = storeCredential.ExpiresAt
	existing.DeletedAt.Valid = false

	return &existing, nil
}

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
