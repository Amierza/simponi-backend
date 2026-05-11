package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IProductVendorRepository interface {
		// CREATE
		CreateProductVendor(ctx context.Context, tx *gorm.DB, productVendor *entity.ProductVendor) error

		// READ

		// UPDATE

		// DELETE
		DeleteProductVendorByProductIDAndVendorID(ctx context.Context, tx *gorm.DB, productID, vendorID *uuid.UUID) error
	}

	productVendorRepository struct {
		db *gorm.DB
	}
)

func NewProductVendorRepository(db *gorm.DB) *productVendorRepository {
	return &productVendorRepository{
		db: db,
	}
}

// CREATE
func (vr *productVendorRepository) CreateProductVendor(ctx context.Context, tx *gorm.DB, productVendor *entity.ProductVendor) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Create(productVendor).Error
}

// READ

// UPDATE

// DELETE
func (vr *productVendorRepository) DeleteProductVendorByProductIDAndVendorID(ctx context.Context, tx *gorm.DB, productID, vendorID *uuid.UUID) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Where("product_id = ?", &productID).Where("vendor_id = ?", vendorID).Delete(&entity.ProductVendor{}).Error
}
