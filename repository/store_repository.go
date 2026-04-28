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
	IStoreRepository interface {
		// CREATE
		CreateStore(ctx context.Context, tx *gorm.DB, store *entity.Store) error

		// READ
		GetStores(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.StorePaginationRepositoryResponse, error)
		GetStoreByID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (*entity.Store, bool, error)

		// UPDATE
		UpdateStore(ctx context.Context, tx *gorm.DB, store *entity.Store) error

		// DELETE
		DeleteStoreByID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) error
	}

	storeRepository struct {
		db *gorm.DB
	}
)

func NewStoreRepository(db *gorm.DB) *storeRepository {
	return &storeRepository{
		db: db,
	}
}

// CREATE
func (vr *storeRepository) CreateStore(ctx context.Context, tx *gorm.DB, store *entity.Store) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Create(store).Error
}

// READ
func (vr *storeRepository) GetStores(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.StorePaginationRepositoryResponse, error) {
	if tx == nil {
		tx = vr.db
	}

	var stores []*entity.Store
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Preload("StorePlatforms").
		Preload("Orders").
		Preload("Logs")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.StorePaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&stores).Error; err != nil {
		return dto.StorePaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.StorePaginationRepositoryResponse{
		Stores: stores,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (vr *storeRepository) GetStoreByID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (*entity.Store, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var store *entity.Store
	err := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Preload("StorePlatforms").
		Preload("Orders").
		Preload("Logs").
		Where("id = ?", storeID).
		First(&store).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return store, true, nil
}

// UPDATE
func (vr *storeRepository) UpdateStore(ctx context.Context, tx *gorm.DB, store *entity.Store) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).
		Model(&entity.Store{}).
		Where("id = ?", store.ID).
		Updates(map[string]interface{}{
			"name":        store.Name,
			"description": store.Description,
			"image_url":   store.ImageURL,
			"is_active":   store.IsActive,
		}).Error
}

// DELETE
func (vr *storeRepository) DeleteStoreByID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &storeID).Delete(&entity.Store{}).Error
}
