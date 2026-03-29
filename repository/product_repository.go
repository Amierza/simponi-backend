package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	IProductRepository interface {
		CreateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error)
		GetProducts(ctx context.Context, tx *gorm.DB) ([]entity.Product, error)
		GetProductByID(ctx context.Context, tx *gorm.DB, productID string) (*entity.Product, error)
		GetProductBySKU(ctx context.Context, tx *gorm.DB, sku string) (*entity.Product, error)
		GetProductsByCategoryID(ctx context.Context, tx *gorm.DB, categoryId string) ([]entity.Product, error)
		UpdateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error)
		DeleteProduct(ctx context.Context, tx *gorm.DB, productID string) error

		UpdateStock(ctx context.Context, tx *gorm.DB, productID string, change int) error
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

func (pr *productRepository) GetProducts(ctx context.Context, tx *gorm.DB) ([]entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	var products []entity.Product

	if err := tx.WithContext(ctx).Model(&entity.Product{}).Preload("Category").Preload("Images").Preload("ExternalProducts").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *productRepository) GetProductByID(ctx context.Context, tx *gorm.DB, productID string) (*entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	var product entity.Product

	err := tx.WithContext(ctx).Model(&entity.Product{}).Preload("Category").Preload("Images").Preload("ExternalProducts").Preload("Logs").Where("id = ?", productID).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (pr *productRepository) GetProductBySKU(ctx context.Context, tx *gorm.DB, sku string) (*entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	var product entity.Product

	err := tx.WithContext(ctx).Model(&entity.Product{}).Preload("Category").Preload("Images").Preload("ExternalProducts").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (pr *productRepository) GetProductsByCategoryID(ctx context.Context, tx *gorm.DB, categoryId string) ([]entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	var products []entity.Product

	if err := tx.WithContext(ctx).Model(&entity.Product{}).Preload("Category").Preload("Images").Preload("ExternalProducts").Where("category_id = ?", categoryId).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, tx *gorm.DB, product *entity.Product) (*entity.Product, error) {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Save(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (pr *productRepository) UpdateStock(ctx context.Context, tx *gorm.DB, productID string, change int) error {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", productID).UpdateColumn("stock", gorm.Expr("stock + ?", change)).Error; err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) DeleteProduct(ctx context.Context, tx *gorm.DB, productID string) error {
	if tx == nil {
		tx = pr.db
	}

	if err := tx.WithContext(ctx).Where("id = ?", productID).Delete(&entity.Product{}).Error; err != nil {
		return err
	}

	return nil
}
