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
		GetProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.ProductPaginationRepositoryResponse, error)
		GetProductStats(ctx context.Context, tx *gorm.DB) (dto.ProductStatsResponse, error)
		GetProductByID(ctx context.Context, tx *gorm.DB, productID *uuid.UUID) (*entity.Product, bool, error)
		GetProductBySKU(ctx context.Context, tx *gorm.DB, sku string) (*entity.Product, bool, error)
		GetProductsByCategoryID(ctx context.Context, tx *gorm.DB, categoryId *uuid.UUID) ([]entity.Product, error)
		UpdateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error)
		DeleteProductByID(ctx context.Context, tx *gorm.DB, productID *uuid.UUID) error

		UpdateStock(ctx context.Context, tx *gorm.DB, productID *uuid.UUID, change int) error
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

func (pr *productRepository) GetProducts(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.ProductPaginationRepositoryResponse, error) {
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
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Platform")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue, searchValue)
	}

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

func (pr *productRepository) GetProductStats(ctx context.Context, tx *gorm.DB) (dto.ProductStatsResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var stats dto.ProductStatsResponse

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Count(&stats.TotalProducts).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}
	stats.TotalSKUs = stats.TotalProducts

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Select("COALESCE(SUM(stock), 0)").
		Scan(&stats.StockUnits).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Where("stock > 0 AND stock <= 10").
		Count(&stats.LowStock).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Where("stock = 0").
		Count(&stats.OutOfStock).Error; err != nil {
		return dto.ProductStatsResponse{}, err
	}

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
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

func (pr *productRepository) GetProductByID(ctx context.Context, tx *gorm.DB, productID *uuid.UUID) (*entity.Product, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var product entity.Product

	err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Platform").
		Preload("Logs").
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

func (pr *productRepository) GetProductBySKU(ctx context.Context, tx *gorm.DB, sku string) (*entity.Product, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var product entity.Product

	err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Platform").
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

func (pr *productRepository) GetProductsByCategoryID(ctx context.Context, tx *gorm.DB, categoryId *uuid.UUID) ([]entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	var products []entity.Product

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Category").
		Preload("Images").
		Preload("ExternalProducts").
		Preload("ExternalProducts.StorePlatform").
		Preload("ExternalProducts.StorePlatform.Platform").
		Where("category_id = ?", categoryId).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).
		Where("id = ?", product.ID).
		Updates(product).
		Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *productRepository) UpdateStock(ctx context.Context, tx *gorm.DB, productID *uuid.UUID, change int) error {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).
		Model(&entity.Product{}).
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", change)).
		Error; err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) DeleteProductByID(ctx context.Context, tx *gorm.DB, productID *uuid.UUID) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).
	Where("id = ?", productID).
	Delete(&entity.Product{}).
	Error
}
