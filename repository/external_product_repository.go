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
	IExternalProductRepository interface {
		CreateExternalProduct(ctx context.Context, tx *gorm.DB, externalProduct *entity.ExternalProduct) (*entity.ExternalProduct, error)
		GetExternalProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ExternalProductPaginationRepositoryResponse, error)
		GetExternalProductByStoreIDAndExprodID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, externalProductID *uuid.UUID) (*entity.ExternalProduct, bool, error)
		GetExternalProductsByStoreIDAndStorePlatformID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, storePlatformID *uuid.UUID) ([]entity.ExternalProduct, error)
		UpdateExternalProductByStoreIDAndExprodID(ctx context.Context, tx *gorm.DB, externalProduct *entity.ExternalProduct) (*entity.ExternalProduct, error)
		DeleteExternalProductByStoreIDAndExprodID(ctx context.Context, tx *gorm.DB, externalProductID *uuid.UUID) error
	}

	externalProductRepository struct {
		db *gorm.DB
	}
)

func NewExternalProductRepository(db *gorm.DB) *externalProductRepository {
	return &externalProductRepository{
		db: db,
	}
}

func (epr *externalProductRepository) CreateExternalProduct(ctx context.Context, tx *gorm.DB, externalProduct *entity.ExternalProduct) (*entity.ExternalProduct, error) {
	if tx == nil {
		tx = epr.db
	}

	if err := tx.WithContext(ctx).Create(externalProduct).Error; err != nil {
		return nil, err
	}

	return externalProduct, nil
}

func (epr *externalProductRepository) GetExternalProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ExternalProductPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = epr.db
	}

	var externalProducts []entity.ExternalProduct
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.ExternalProduct{}).
		Joins("JOIN store_platforms sp ON sp.id = external_products.store_platform_id").
		Where("sp.store_id = ?", storeID).
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
		Preload("StorePlatform.Store").
		Preload("StorePlatform.Platform")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.
			Joins("LEFT JOIN products ON products.id = external_products.product_id").
			Where("LOWER(products.name) LIKE ? OR LOWER(products.sku) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.ExternalProductPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"external_products"."created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&externalProducts).Error; err != nil {
		return dto.ExternalProductPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.ExternalProductPaginationRepositoryResponse{
		ExternalProducts: externalProducts,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (epr *externalProductRepository) GetExternalProductByStoreIDAndExprodID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, externalProductID *uuid.UUID) (*entity.ExternalProduct, bool, error) {
	if tx == nil {
		tx = epr.db
	}

	var externalProduct entity.ExternalProduct

	err := tx.WithContext(ctx).
		Model(&entity.ExternalProduct{}).
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
		Preload("StorePlatform.Store").
		Preload("StorePlatform.Platform").
		Joins("JOIN store_platforms sp ON sp.id = external_products.store_platform_id").
		Where("sp.store_id = ?", storeID).
		Where("external_products.id = ?", externalProductID).
		First(&externalProduct).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &externalProduct, true, nil
}

func (epr *externalProductRepository) GetExternalProductsByStoreIDAndStorePlatformID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, storePlatformID *uuid.UUID) ([]entity.ExternalProduct, error) {
	if tx == nil {
		tx = epr.db
	}

	var externalProducts []entity.ExternalProduct

	if err := tx.WithContext(ctx).
		Model(&entity.ExternalProduct{}).
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
		Preload("StorePlatform.Store").
		Preload("StorePlatform.Platform").
		Joins("JOIN store_platforms sp ON sp.ID = external_products.store_platform_id").
		Where("sp.id = ?", storePlatformID).
		Where("sp.store_id = ?", storeID).
		Find(&externalProducts).Error; err != nil {
		return nil, err
	}

	return externalProducts, nil
}

func (epr *externalProductRepository) UpdateExternalProductByStoreIDAndExprodID(ctx context.Context, tx *gorm.DB, externalProduct *entity.ExternalProduct) (*entity.ExternalProduct, error) {
	if tx == nil {
		tx = epr.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", externalProduct.ID).Updates(externalProduct).Error; err != nil {
		return nil, err
	}

	return externalProduct, nil
}

func (epr *externalProductRepository) DeleteExternalProductByStoreIDAndExprodID(ctx context.Context, tx *gorm.DB, externalProductID *uuid.UUID) error {
	if tx == nil {
		tx = epr.db
	}

	return tx.WithContext(ctx).Where("id = ?", externalProductID).Delete(&entity.ExternalProduct{}).Error
}
