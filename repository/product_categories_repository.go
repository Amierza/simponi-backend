package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	IProductCategoriesRepository interface {
		GetProductCategories(ctx context.Context, tx *gorm.DB) ([]entity.ProductCategory, error)
	}

	productCategoriesRepository struct {
		db *gorm.DB
	}
)

func NewProductCategoriesRepository(db *gorm.DB) *productCategoriesRepository {
	return &productCategoriesRepository{
		db: db,
	}
}


func (pr *productCategoriesRepository) GetProductCategories(ctx context.Context, tx *gorm.DB) ([]entity.ProductCategory, error) {
	if tx == nil {
		tx = pr.db
	}

	var categories []entity.ProductCategory

	if err := tx.WithContext(ctx).
		Model(&entity.ProductCategory{}).
		Find(&categories).Error; err != nil {
			return nil, err
		}
	return categories, nil
}