package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IPlatformRepository interface {
		// READ
		GetPlatformByPlatformID(ctx context.Context, tx *gorm.DB, platformID *uuid.UUID) (*entity.Platform, bool, error)
	}

	platformRepository struct {
		db *gorm.DB
	}
)

func NewPlatformRepository(db *gorm.DB) *platformRepository {
	return &platformRepository{
		db: db,
	}
}

// READ
func (pr *platformRepository) GetPlatformByPlatformID(ctx context.Context, tx *gorm.DB, platformID *uuid.UUID) (*entity.Platform, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	platform := new(entity.Platform)
	err := tx.WithContext(ctx).
		Preload("StorePlatforms.Store").
		Where("id = ?", platformID).
		Take(platform).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return platform, true, nil
}
