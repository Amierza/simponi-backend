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
	IVendorRepository interface {
		// CREATE
		CreateVendor(ctx context.Context, tx *gorm.DB, vendor *entity.Vendor) error

		// READ
		GetVendorsByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, req *response.PaginationRequest) (dto.VendorPaginationRepositoryResponse, error)
		GetVendorByStoreIDAndVendorID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, vendorID *uuid.UUID) (*entity.Vendor, bool, error)
		GetVendorByStoreIDAndVendorPhoneNumber(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, phoneNumber string) (*entity.Vendor, bool, error)
		GetVendorByStoreIDAndVendorEmail(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, email string) (*entity.Vendor, bool, error)

		// UPDATE
		UpdateVendor(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, vendor *entity.Vendor) error

		// DELETE
		DeleteVendorByStoreIDAndVendorID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, vendorID *uuid.UUID) error
	}

	vendorRepository struct {
		db *gorm.DB
	}
)

func NewVendorRepository(db *gorm.DB) *vendorRepository {
	return &vendorRepository{
		db: db,
	}
}

// CREATE
func (vr *vendorRepository) CreateVendor(ctx context.Context, tx *gorm.DB, vendor *entity.Vendor) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Create(vendor).Error
}

// READ
func (vr *vendorRepository) GetVendorsByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, req *response.PaginationRequest) (dto.VendorPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendors []*entity.Vendor
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Where("store_id = ?", storeID).
		Preload("Store").
		Preload("ProductVendors")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone_number) LIKE ? OR LOWER(address) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue, searchValue, searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.VendorPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&vendors).Error; err != nil {
		return dto.VendorPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.VendorPaginationRepositoryResponse{
		Vendors: vendors,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (vr *vendorRepository) GetVendorByStoreIDAndVendorID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, vendorID *uuid.UUID) (*entity.Vendor, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor *entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("Store").
		Preload("ProductVendors").
		Where("id = ?", vendorID).
		Where("store_id = ?", storeID).
		First(&vendor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return vendor, true, nil
}
func (vr *vendorRepository) GetVendorByStoreIDAndVendorPhoneNumber(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, phoneNumber string) (*entity.Vendor, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("Store").
		Preload("ProductVendors").
		Where("phone_number = ?", phoneNumber).
		Where("store_id = ?", storeID).
		First(&vendor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &vendor, true, nil
}
func (vr *vendorRepository) GetVendorByStoreIDAndVendorEmail(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, email string) (*entity.Vendor, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("Store").
		Preload("ProductVendors").
		Where("email = ?", email).
		Where("store_id = ?", storeID).
		First(&vendor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &vendor, true, nil
}

// UPDATE
func (vr *vendorRepository) UpdateVendor(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, vendor *entity.Vendor) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Model(&entity.Vendor{}).Where("id = ?", vendor.ID).Where("store_id = ?", storeID).Updates(&vendor).Error
}

// DELETE
func (vr *vendorRepository) DeleteVendorByStoreIDAndVendorID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, vendorID *uuid.UUID) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &vendorID).Where("store_id = ?", storeID).Delete(&entity.Vendor{}).Error
}
