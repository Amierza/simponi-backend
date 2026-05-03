package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IStorePlatformRepository interface {
		// CREATE
		CreateStorePlatform(ctx context.Context, tx *gorm.DB, storePlatform *entity.StorePlatform) error

		// READ
		GetStorePlatformByID(ctx context.Context, tx *gorm.DB, id *uuid.UUID) (*entity.StorePlatform, bool, error)
		GetStorePlatformsByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) ([]*entity.StorePlatform, error)
		GetStorePlatformByStoreIDAndPlatformID(ctx context.Context, tx *gorm.DB, storeID, platformID *uuid.UUID) (*entity.StorePlatform, bool, error)
		CountStorePlatformsByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (int64, error)

		// DELETE
		DeleteStorePlatformByID(ctx context.Context, tx *gorm.DB, id *uuid.UUID) error
	}

	storePlatformRepository struct {
		db *gorm.DB
	}
)

func NewStorePlatformRepository(db *gorm.DB) *storePlatformRepository {
	return &storePlatformRepository{
		db: db,
	}
}

// CREATE
func (spr *storePlatformRepository) CreateStorePlatform(ctx context.Context, tx *gorm.DB, storePlatform *entity.StorePlatform) error {
	if tx == nil {
		tx = spr.db
	}

	return tx.WithContext(ctx).Create(storePlatform).Error
}

// READ
func (r *storePlatformRepository) GetStorePlatformByID(
	ctx context.Context, tx *gorm.DB, id *uuid.UUID,
) (*entity.StorePlatform, bool, error) {
	if tx == nil {
		tx = r.db
	}
	var sp entity.StorePlatform
	err := tx.WithContext(ctx).
		Preload("Platform").
		Preload("Credential").
		Where("id = ?", id).
		First(&sp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &sp, true, nil
}

func (r *storePlatformRepository) GetStorePlatformsByStoreID(
	ctx context.Context, tx *gorm.DB, storeID *uuid.UUID,
) ([]*entity.StorePlatform, error) {
	if tx == nil {
		tx = r.db
	}
	var sps []*entity.StorePlatform
	err := tx.WithContext(ctx).
		Preload("Platform").
		Where("store_id = ?", storeID).
		Find(&sps).Error
	return sps, err
}

func (r *storePlatformRepository) GetStorePlatformByStoreIDAndPlatformID(
	ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, platformID *uuid.UUID,
) (*entity.StorePlatform, bool, error) {
	if tx == nil {
		tx = r.db
	}
	var sp entity.StorePlatform
	err := tx.WithContext(ctx).
		Where("store_id = ? AND platform_id = ?", storeID, platformID).
		First(&sp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &sp, true, nil
}
func (r *storePlatformRepository) CountStorePlatformsByStoreID(
	ctx context.Context, tx *gorm.DB, storeID *uuid.UUID,
) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var count int64
	err := tx.WithContext(ctx).
		Model(&entity.StorePlatform{}).
		Where("store_id = ?", storeID).
		Count(&count).Error
	return count, err
}

// DELETE
func (r *storePlatformRepository) DeleteStorePlatformByID(
	ctx context.Context, tx *gorm.DB, id *uuid.UUID,
) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.StorePlatform{}).Error
}
