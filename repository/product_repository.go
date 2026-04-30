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
	IProductRepository interface {
		CreateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error)
		CreateProductImage(ctx context.Context, tx *gorm.DB, productImage *entity.ProductImage) (*entity.ProductImage, error)
		AttachProductImageToProduct(ctx context.Context, tx *gorm.DB, imageID *uuid.UUID, productID *uuid.UUID) error
		GetProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ProductPaginationRepositoryResponse, error)
		GetProductStats(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (dto.ProductStatsResponse, error)
		GetProductByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, productID *uuid.UUID) (*entity.Product, bool, error)
		GetProductBySKUAndStoreID(ctx context.Context, tx *gorm.DB, sku string, storeID *uuid.UUID) (*entity.Product, bool, error)
		UpdateProductByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error)
		UpdateStockByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, productID *uuid.UUID, change int) error
		DeleteProductByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, productID *uuid.UUID) error
	}

	productRepository struct {
		db *gorm.DB
	}
)

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (pr *productRepository) CreateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *productRepository) CreateProductImage(ctx context.Context, tx *gorm.DB, productImage *entity.ProductImage) (*entity.ProductImage, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Create(productImage).Error; err != nil {
		return nil, err
	}

	return productImage, nil
}

func (pr *productRepository) AttachProductImageToProduct(ctx context.Context, tx *gorm.DB, imageID *uuid.UUID, productID *uuid.UUID) error {
	if tx == nil {
		tx = pr.db
	}

	result := tx.WithContext(ctx).
		Model(&entity.ProductImage{}).
		Where("id = ?", imageID).
		Where("product_id IS NULL").
		Update("product_id", productID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (pr *productRepository) GetProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ProductPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var products []entity.Product
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Where("store_id = ?", storeID).
		Preload("Store").
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Store").
		Preload("ExternalProducts.StorePlatform.Platform")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue, searchValue)
	}

	// if req.SKU != "" {
	// 	query = query.Where("sku = ?", req.SKU)
	// }

	// if req.CategoryID != "" {
	// 	query = query.Where("category_id = ?", req.CategoryID)
	// }

	if err := query.Count(&count).Error; err != nil {
		return dto.ProductPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&products).Error; err != nil {
		return dto.ProductPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.ProductPaginationRepositoryResponse{
		Products: products,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (pr *productRepository) GetProductStats(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (dto.ProductStatsResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var stats dto.ProductStatsResponse

	baseQuery := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Where("store_id = ?", *storeID)

	if err := baseQuery.
		Count(&stats.TotalProducts).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}
	stats.TotalSKUs = stats.TotalProducts

	if err := baseQuery.
		Select("COALESCE(SUM(stock), 0)").
		Scan(&stats.StockUnits).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	if err := baseQuery.
		Where("stock > 0 AND stock <= 10").
		Count(&stats.LowStock).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	if err := baseQuery.
		Where("stock = 0").
		Count(&stats.OutOfStock).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	if err := baseQuery.
		Where("id NOT IN (?)",
			tx.Model(&entity.ExternalProduct{}).
				Select("DISTINCT product_id").
				Where("product_id IS NOT NULL"),
		).
		Count(&stats.Unsynced).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	return stats, nil
}

func (pr *productRepository) GetProductByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, productID *uuid.UUID) (*entity.Product, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var product entity.Product

	err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Store").
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Store").
		Preload("ExternalProducts.StorePlatform.Platform").
		Preload("Logs").
		Where("store_id = ?", storeID).
		Where("id = ?", productID).
		First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &product, true, nil
}

func (pr *productRepository) GetProductBySKUAndStoreID(ctx context.Context, tx *gorm.DB, sku string, storeID *uuid.UUID) (*entity.Product, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var product entity.Product

	err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Store").
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Store").
		Preload("ExternalProducts.StorePlatform.Platform").
		Where("store_id = ?", storeID).
		Where("sku = ?", sku).
		First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &product, true, nil
}

func (pr *productRepository) UpdateProductByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).
		Where("store_id = ?", product.StoreID).
		Where("id = ?", product.ID).
		Updates(product).
		Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *productRepository) UpdateStockByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, productID *uuid.UUID, change int) error {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Where("store_id = ?", productID).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", change)).
		Error; err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) DeleteProductByStoreIDAndProductID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, productID *uuid.UUID) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).
		Where("store_id = ?", storeID).
		Where("id = ?", productID).
		Delete(&entity.Product{}).
		Error
}
