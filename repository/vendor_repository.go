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
		GetVendors(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.VendorPaginationRepositoryResponse, error)
		GetVendorByID(ctx context.Context, tx *gorm.DB, vendorID *uuid.UUID) (*entity.Vendor, bool, error)
		// ✅ FIX: tambah GetVendorByName ke interface
		GetVendorByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Vendor, error)
		GetVendorByPhoneNumber(ctx context.Context, tx *gorm.DB, phoneNumber string) (*entity.Vendor, bool, error)
		GetVendorByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.Vendor, bool, error)

		// UPDATE
		UpdateVendor(ctx context.Context, tx *gorm.DB, vendor *entity.Vendor) error

		// DELETE
		DeleteVendorByID(ctx context.Context, tx *gorm.DB, vendorID *uuid.UUID) error
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
func (vr *vendorRepository) GetVendors(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.VendorPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendors []*entity.Vendor
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("ProductVendors")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(
			"LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone_number) LIKE ? OR LOWER(address) LIKE ? OR LOWER(description) LIKE ?",
			searchValue, searchValue, searchValue, searchValue, searchValue,
		)
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
	}, nil
}

func (vr *vendorRepository) GetVendorByID(ctx context.Context, tx *gorm.DB, vendorID *uuid.UUID) (*entity.Vendor, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor *entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("ProductVendors").
		Where("id = ?", vendorID).
		First(&vendor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return vendor, true, nil
}

// ✅ FIX: implementasi GetVendorByName — case-insensitive
// Return nil jika tidak ditemukan (bukan error), return *entity.Vendor jika ada
func (vr *vendorRepository) GetVendorByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Vendor, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Where("LOWER(name) = LOWER(?)", name).
		First(&vendor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // tidak ditemukan = bukan error
	}
	if err != nil {
		return nil, err
	}

	return &vendor, nil
}

func (vr *vendorRepository) GetVendorByPhoneNumber(ctx context.Context, tx *gorm.DB, phoneNumber string) (*entity.Vendor, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("ProductVendors").
		Where("phone_number = ?", phoneNumber).
		First(&vendor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &vendor, true, nil
}

func (vr *vendorRepository) GetVendorByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.Vendor, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var vendor entity.Vendor
	err := tx.WithContext(ctx).
		Model(&entity.Vendor{}).
		Preload("ProductVendors").
		Where("email = ?", email).
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
func (vr *vendorRepository) UpdateVendor(ctx context.Context, tx *gorm.DB, vendor *entity.Vendor) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Model(&entity.Vendor{}).Where("id = ?", vendor.ID).Updates(&vendor).Error
}

// DELETE
func (vr *vendorRepository) DeleteVendorByID(ctx context.Context, tx *gorm.DB, vendorID *uuid.UUID) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &vendorID).Delete(&entity.Vendor{}).Error
}
