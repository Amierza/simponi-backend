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
		GetStorePlatformByStoreIDAndPlatformID(ctx context.Context, tx *gorm.DB, storeID, platformID *uuid.UUID) (*entity.StorePlatform, bool, error)
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
func (spr *storePlatformRepository) GetStorePlatformByStoreIDAndPlatformID(ctx context.Context, tx *gorm.DB, storeID, platformID *uuid.UUID) (*entity.StorePlatform, bool, error) {
	if tx == nil {
		tx = spr.db
	}

	storePlatform := new(entity.StorePlatform)
	err := tx.WithContext(ctx).
		Preload("Store").
		Preload("Platform").
		Where("store_id = ?", storeID).
		Where("platform_id = ?", platformID).
		First(storePlatform).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return storePlatform, true, nil
}
