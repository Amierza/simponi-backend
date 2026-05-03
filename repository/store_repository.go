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
		GetStoresByUserID(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, userID *uuid.UUID) (dto.StorePaginationRepositoryResponse, error)
		GetStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (*entity.Store, bool, error)
		GetStoreByUserID(ctx context.Context, tx *gorm.DB, userID *uuid.UUID) (*entity.Store, bool, error)

		// UPDATE
		UpdateStoreByStoreID(ctx context.Context, tx *gorm.DB, store *entity.Store) error

		// DELETE
		DeleteStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) error
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
		Preload("StorePlatforms.Platform").
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
func (vr *storeRepository) GetStoresByUserID(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, userID *uuid.UUID) (dto.StorePaginationRepositoryResponse, error) {
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
		Joins("JOIN store_users su ON su.store_id = stores.id").
		Where("su.user_id = ?", userID).
		Preload("StorePlatforms.Platform").
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
func (vr *storeRepository) GetStoreByUserID(
	ctx context.Context, tx *gorm.DB, userID *uuid.UUID,
) (*entity.Store, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	// Scan ke string dulu — GORM tidak bisa auto-convert UUID string ke uuid.UUID
	var storeIDStr string
	err := tx.WithContext(ctx).
		Table("store_users").
		Select("store_id").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(1).
		Scan(&storeIDStr).Error
	if err != nil {
		return nil, false, err
	}

	// Scan tidak return ErrRecordNotFound, cek manual
	if storeIDStr == "" {
		return nil, false, nil
	}

	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, false, err
	}

	var store entity.Store
	err = tx.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", storeID).
		Preload("StorePlatforms", "deleted_at IS NULL").
		Preload("StorePlatforms.Platform").
		First(&store).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &store, true, nil
}

func (vr *storeRepository) GetStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (*entity.Store, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var store *entity.Store
	err := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Preload("StorePlatforms.Platform").
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
func (vr *storeRepository) UpdateStoreByStoreID(ctx context.Context, tx *gorm.DB, store *entity.Store) error {
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
func (vr *storeRepository) DeleteStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &storeID).Delete(&entity.Store{}).Error
}
