package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/helper"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IVendorService interface {
		CreateVendor(ctx context.Context, req *dto.CreateVendorRequest) (*dto.VendorResponse, error)
		GetVendors(ctx context.Context, req *response.PaginationRequest) (dto.VendorPaginationResponse, error)
		GetVendorByID(ctx context.Context, vendorID *uuid.UUID) (*dto.VendorResponse, error)
		UpdateVendor(ctx context.Context, vendorID *uuid.UUID, req *dto.UpdateVendorRequest) (*dto.VendorResponse, error)
		DeleteVendorByID(ctx context.Context, vendorID *uuid.UUID) error
	}

	vendorService struct {
		vendorRepo repository.IVendorRepository
		logger     *zap.Logger
		jwtService jwt.IJWT
	}
)

func NewVendorService(vendorRepo repository.IVendorRepository, logger *zap.Logger, jwtService jwt.IJWT) *vendorService {
	return &vendorService{
		vendorRepo: vendorRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}

func mapToVendorResponse(v *entity.Vendor) *dto.VendorResponse {
	return &dto.VendorResponse{
		ID:          v.ID,
		Name:        v.Name,
		Email:       v.Email,
		PhoneNumber: v.PhoneNumber,
		Address:     v.Address,
		ImageURL:    v.ImageURL,
		Description: v.Description,
	}
}

func (vs *vendorService) CreateVendor(ctx context.Context, req *dto.CreateVendorRequest) (*dto.VendorResponse, error) {
	// ✅ FIX: check duplicate name
	existingByName, err := vs.vendorRepo.GetVendorByName(ctx, nil, req.Name)
	if err != nil {
		vs.logger.Error("failed to check vendor name", zap.String("name", req.Name), zap.Error(err))
		return nil, fmt.Errorf("failed to check vendor name: %w", dto.ErrInternal)
	}
	if existingByName != nil {
		vs.logger.Warn("vendor name already exists", zap.String("name", req.Name))
		return nil, fmt.Errorf("vendor with name '%s' already exists: %w", req.Name, dto.ErrAlreadyExists)
	}

	// check duplicate email
	if req.Email != "" {
		_, found, err := vs.vendorRepo.GetVendorByEmail(ctx, nil, req.Email)
		if err != nil {
			vs.logger.Error("failed to get vendor by email", zap.String("email", req.Email), zap.Error(err))
			return nil, fmt.Errorf("failed to check vendor email: %w", dto.ErrInternal)
		}
		if found {
			vs.logger.Warn("vendor email already exists", zap.String("email", req.Email))
			return nil, fmt.Errorf("vendor with email '%s' already exists: %w", req.Email, dto.ErrAlreadyExists)
		}
	}

	// validate & normalize phone number
	phoneNumber, err := helper.NormalizePhoneNumber(req.PhoneNumber)
	if err != nil {
		vs.logger.Error("invalid phone number", zap.String("phone_number", req.PhoneNumber), zap.Error(err))
		return nil, fmt.Errorf("invalid phone number: %w", dto.ErrBadRequest)
	}

	_, found, err := vs.vendorRepo.GetVendorByPhoneNumber(ctx, nil, phoneNumber)
	if err != nil {
		vs.logger.Error("failed to get vendor by phone number", zap.String("phone_number", phoneNumber), zap.Error(err))
		return nil, fmt.Errorf("failed to check vendor phone number: %w", dto.ErrInternal)
	}
	if found {
		vs.logger.Warn("vendor phone number already exists", zap.String("phone_number", phoneNumber))
		return nil, fmt.Errorf("vendor with phone number '%s' already exists: %w", phoneNumber, dto.ErrAlreadyExists)
	}

	newID := uuid.New()
	newVendor := &entity.Vendor{
		ID:          newID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: phoneNumber,
		Address:     req.Address,
		ImageURL:    req.ImageURL,
		Description: req.Description,
	}

	err = vs.vendorRepo.CreateVendor(ctx, nil, newVendor)
	if err != nil {
		vs.logger.Error("failed to create vendor", zap.Error(err))
		return nil, fmt.Errorf("failed to create vendor: %w", dto.ErrInternal)
	}

	vs.logger.Info("success to create vendor", zap.String("id", newVendor.ID.String()))

	return mapToVendorResponse(newVendor), nil
}

func (vs *vendorService) GetVendors(ctx context.Context, req *response.PaginationRequest) (dto.VendorPaginationResponse, error) {
	datas, err := vs.vendorRepo.GetVendors(ctx, nil, req)
	if err != nil {
		vs.logger.Error("failed to get vendors", zap.Error(err))
		return dto.VendorPaginationResponse{}, fmt.Errorf("failed to get vendors: %w", dto.ErrInternal)
	}

	vs.logger.Info("success to get vendors", zap.Int64("count", datas.Count))

	var vendors []*dto.VendorResponse
	for _, vendor := range datas.Vendors {
		vendors = append(vendors, mapToVendorResponse(vendor))
	}

	return dto.VendorPaginationResponse{
		Data:               vendors,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (vs *vendorService) GetVendorByID(ctx context.Context, vendorID *uuid.UUID) (*dto.VendorResponse, error) {
	vendor, found, err := vs.vendorRepo.GetVendorByID(ctx, nil, vendorID)
	if err != nil {
		vs.logger.Error("failed to get vendor by ID", zap.String("vendorID", vendorID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get vendor ID: %w", dto.ErrInternal)
	}
	if !found {
		vs.logger.Warn("vendor not found", zap.String("vendorID", vendorID.String()))
		return nil, fmt.Errorf("vendor not found: %v", dto.ErrNotFound)
	}

	vs.logger.Info("success to get vendor by id", zap.String("id", vendorID.String()))

	return mapToVendorResponse(vendor), nil
}

func (vs *vendorService) UpdateVendor(ctx context.Context, vendorID *uuid.UUID, req *dto.UpdateVendorRequest) (*dto.VendorResponse, error) {
	vendor, found, err := vs.vendorRepo.GetVendorByID(ctx, nil, vendorID)
	if err != nil {
		vs.logger.Error("failed to get vendor by ID", zap.String("vendorID", vendorID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get vendor ID: %w", dto.ErrInternal)
	}
	if !found {
		vs.logger.Warn("vendor not found", zap.String("vendorID", vendorID.String()))
		return nil, fmt.Errorf("vendor not found: %v", dto.ErrNotFound)
	}

	// ✅ FIX: check duplicate name, skip jika nama tidak berubah
	if !strings.EqualFold(vendor.Name, req.Name) {
		existingByName, err := vs.vendorRepo.GetVendorByName(ctx, nil, req.Name)
		if err != nil {
			vs.logger.Error("failed to check vendor name", zap.String("name", req.Name), zap.Error(err))
			return nil, fmt.Errorf("failed to check vendor name: %w", dto.ErrInternal)
		}
		if existingByName != nil {
			vs.logger.Warn("vendor name already exists", zap.String("name", req.Name))
			return nil, fmt.Errorf("vendor with name '%s' already exists: %w", req.Name, dto.ErrAlreadyExists)
		}
	}
	vendor.Name = req.Name

	// validate email — skip jika email tidak berubah
	if req.Email != nil {
		if vendor.Email != *req.Email {
			_, found, err = vs.vendorRepo.GetVendorByEmail(ctx, nil, *req.Email)
			if err != nil {
				vs.logger.Error("failed to get vendor by email", zap.String("email", *req.Email), zap.Error(err))
				return nil, fmt.Errorf("failed to check vendor email: %w", dto.ErrInternal)
			}
			if found {
				vs.logger.Warn("vendor email already exists", zap.String("email", *req.Email))
				return nil, fmt.Errorf("vendor with email '%s' already exists: %w", *req.Email, dto.ErrAlreadyExists)
			}
		}
		vendor.Email = *req.Email
	}

	// validate & normalize phone number — skip jika phone tidak berubah
	phoneNumber, err := helper.NormalizePhoneNumber(req.PhoneNumber)
	if err != nil {
		vs.logger.Error("invalid phone number", zap.String("phone_number", req.PhoneNumber), zap.Error(err))
		return nil, fmt.Errorf("invalid phone number: %w", dto.ErrBadRequest)
	}
	if vendor.PhoneNumber != phoneNumber {
		_, found, err = vs.vendorRepo.GetVendorByPhoneNumber(ctx, nil, phoneNumber)
		if err != nil {
			vs.logger.Error("failed to get vendor by phone number", zap.String("phone_number", phoneNumber), zap.Error(err))
			return nil, fmt.Errorf("failed to check vendor phone number: %w", dto.ErrInternal)
		}
		if found {
			vs.logger.Warn("vendor phone number already exists", zap.String("phone_number", phoneNumber))
			return nil, fmt.Errorf("vendor with phone number '%s' already exists: %w", phoneNumber, dto.ErrAlreadyExists)
		}
		vendor.PhoneNumber = phoneNumber
	}

	if req.Address != nil {
		vendor.Address = *req.Address
	}
	if req.ImageURL != nil {
		vendor.ImageURL = *req.ImageURL
	}
	if req.Description != nil {
		vendor.Description = *req.Description
	}

	err = vs.vendorRepo.UpdateVendor(ctx, nil, vendor)
	if err != nil {
		vs.logger.Error("failed to update vendor", zap.String("id", vendorID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update vendor: %w", dto.ErrInternal)
	}

	return mapToVendorResponse(vendor), nil
}

func (vs *vendorService) DeleteVendorByID(ctx context.Context, vendorID *uuid.UUID) error {
	_, found, err := vs.vendorRepo.GetVendorByID(ctx, nil, vendorID)
	if err != nil {
		vs.logger.Error("failed to get vendor by ID", zap.String("vendorID", vendorID.String()), zap.Error(err))
		return fmt.Errorf("failed to get vendor ID: %w", dto.ErrInternal)
	}
	if !found {
		vs.logger.Warn("vendor not found", zap.String("vendorID", vendorID.String()))
		return fmt.Errorf("vendor not found: %v", dto.ErrNotFound)
	}

	if err := vs.vendorRepo.DeleteVendorByID(ctx, nil, vendorID); err != nil {
		vs.logger.Error("failed to delete vendor by id", zap.String("vendorID", vendorID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete vendor by id: %w", dto.ErrInternal)
	}

	return nil
}
