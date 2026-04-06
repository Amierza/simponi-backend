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
		GetExternalProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.ExternalProductPaginationRepositoryResponse, error)
		GetExternalProductByID(ctx context.Context, tx *gorm.DB, externalProductID *uuid.UUID) (*entity.ExternalProduct, bool, error)
		GetExternalProductByProductID(ctx context.Context, tx *gorm.DB, productID *uuid.UUID) ([]entity.ExternalProduct, error)
		GetExternalProductByStorePlatformID(ctx context.Context, tx *gorm.DB, storePlatformID *uuid.UUID) ([]entity.ExternalProduct, error)
		UpdateExternalProduct(ctx context.Context, tx *gorm.DB, externalProduct *entity.ExternalProduct) (*entity.ExternalProduct, error)
		DeleteExternalProductByID(ctx context.Context, tx *gorm.DB, externalProductID *uuid.UUID) error
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

func (epr *externalProductRepository) GetExternalProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.ExternalProductPaginationRepositoryResponse, error) {
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
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
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

func (epr *externalProductRepository) GetExternalProductByID(ctx context.Context, tx *gorm.DB, externalProductID *uuid.UUID) (*entity.ExternalProduct, bool, error) {
	if tx == nil {
		tx = epr.db
	}

	var externalProduct entity.ExternalProduct

	err := tx.WithContext(ctx).
		Model(&entity.ExternalProduct{}).
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
		Preload("StorePlatform.Platform").
		Where("id = ?", externalProductID).
		First(&externalProduct).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &externalProduct, true, nil
}

func (epr *externalProductRepository) GetExternalProductByProductID(ctx context.Context, tx *gorm.DB, productID *uuid.UUID) ([]entity.ExternalProduct, error) {
	if tx == nil {
		tx = epr.db
	}

	var externalProducts []entity.ExternalProduct

	if err := tx.WithContext(ctx).
		Model(&entity.ExternalProduct{}).
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
		Preload("StorePlatform.Platform").
		Where("product_id = ?", productID).
		Find(&externalProducts).Error; err != nil {
		return nil, err
	}

	return externalProducts, nil
}

func (epr *externalProductRepository) GetExternalProductByStorePlatformID(ctx context.Context, tx *gorm.DB, storePlatformID *uuid.UUID) ([]entity.ExternalProduct, error) {
	if tx == nil {
		tx = epr.db
	}

	var externalProducts []entity.ExternalProduct

	if err := tx.WithContext(ctx).
		Model(&entity.ExternalProduct{}).
		Preload("Product").
		Preload("Product.Images").
		Preload("StorePlatform").
		Preload("StorePlatform.Platform").
		Where("store_platform_id = ?", storePlatformID).
		Find(&externalProducts).Error; err != nil {
		return nil, err
	}

	return externalProducts, nil
}

func (epr *externalProductRepository) UpdateExternalProduct(ctx context.Context, tx *gorm.DB, externalProduct *entity.ExternalProduct) (*entity.ExternalProduct, error) {
	if tx == nil {
		tx = epr.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", externalProduct.ID).Updates(externalProduct).Error; err != nil {
		return nil, err
	}

	return externalProduct, nil
}

func (epr *externalProductRepository) DeleteExternalProductByID(ctx context.Context, tx *gorm.DB, externalProductID *uuid.UUID) error {
	if tx == nil {
		tx = epr.db
	}

	return tx.WithContext(ctx).Where("id = ?", externalProductID).Delete(&entity.ExternalProduct{}).Error
}