package service

import (
	"context"
	"fmt"

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
		CreateVendorByStoreID(ctx context.Context, req *dto.CreateVendorRequest) (*dto.VendorResponse, error)
		GetVendorsByStoreID(ctx context.Context, storeID *uuid.UUID, req *response.PaginationRequest) (dto.VendorPaginationResponse, error)
		GetVendorByStoreIDAndVendorID(ctx context.Context, storeID *uuid.UUID, vendorID *uuid.UUID) (*dto.VendorResponse, error)
		UpdateVendorByStoreIDAndVendorID(ctx context.Context, req *dto.UpdateVendorRequest) (*dto.VendorResponse, error)
		DeleteVendorByStoreIDAndVendorID(ctx context.Context, storeID *uuid.UUID, vendorID *uuid.UUID) error
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

func (vs *vendorService) CreateVendorByStoreID(ctx context.Context, req *dto.CreateVendorRequest) (*dto.VendorResponse, error) {
	// check if email already exists
	if req.Email != "" {
		_, found, err := vs.vendorRepo.GetVendorByStoreIDAndVendorEmail(ctx, nil, req.StoreID, req.Email)
		if err != nil {
			vs.logger.Error("failed to get vendor by store ID and email", zap.String("store_id", req.StoreID.String()), zap.String("email", req.Email), zap.Error(err))
			return nil, fmt.Errorf("failed to get vendor by store ID and email: %w", dto.ErrInternal)
		}
		if found {
			vs.logger.Warn("vendor email already exists", zap.String("email", req.Email))
			return nil, fmt.Errorf("vendor already exists: %w", dto.ErrAlreadyExists)
		}
	}

	// validate & normalize phone number
	phoneNumber, err := helper.NormalizePhoneNumber(req.PhoneNumber)
	if err != nil {
		vs.logger.Error("invalid phone number", zap.String("phone_number", req.PhoneNumber), zap.Error(err))
		return nil, fmt.Errorf("invalid phone number: %w", dto.ErrBadRequest)
	}

	_, found, err := vs.vendorRepo.GetVendorByStoreIDAndVendorPhoneNumber(ctx, nil, req.StoreID, phoneNumber)
	if err != nil {
		vs.logger.Error("failed to get vendor by store ID and phone number", zap.String("store_id", req.StoreID.String()), zap.String("phone_number", phoneNumber), zap.Error(err))
		return nil, fmt.Errorf("failed to get vendor by store ID and phone number: %w", dto.ErrInternal)
	}
	if found {
		vs.logger.Warn("vendor already exists", zap.String("phone_number", phoneNumber))
		return nil, fmt.Errorf("vendor already exists: %w", dto.ErrAlreadyExists)
	}

	newID := uuid.New()
	newVendor := &entity.Vendor{
		ID:          newID,
		StoreID:     *req.StoreID,
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

func (vs *vendorService) GetVendorsByStoreID(ctx context.Context, storeID *uuid.UUID, req *response.PaginationRequest) (dto.VendorPaginationResponse, error) {
	datas, err := vs.vendorRepo.GetVendorsByStoreID(ctx, nil, storeID, req)
	if err != nil {
		vs.logger.Error("failed to get vendors by store ID", zap.Error(err))
		return dto.VendorPaginationResponse{}, fmt.Errorf("failed to get vendors by store ID: %w", dto.ErrInternal)
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

func (vs *vendorService) GetVendorByStoreIDAndVendorID(ctx context.Context, storeID, vendorID *uuid.UUID) (*dto.VendorResponse, error) {
	vendor, found, err := vs.vendorRepo.GetVendorByStoreIDAndVendorID(ctx, nil, storeID, vendorID)
	if err != nil {
		vs.logger.Error("failed to get vendor by store ID and Vendor ID", zap.String("store_id", storeID.String()), zap.String("vendorID", vendorID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get vendor by store ID and Vendor ID: %w", dto.ErrInternal)
	}
	if !found {
		vs.logger.Warn("vendor not found", zap.String("store_id", storeID.String()), zap.String("vendorID", vendorID.String()))
		return nil, fmt.Errorf("vendor not found: %v", dto.ErrNotFound)
	}

	vs.logger.Info("success to get vendor by id", zap.String("id", vendorID.String()))

	return mapToVendorResponse(vendor), nil
}

func (vs *vendorService) UpdateVendorByStoreIDAndVendorID(ctx context.Context, req *dto.UpdateVendorRequest) (*dto.VendorResponse, error) {
	vendor, found, err := vs.vendorRepo.GetVendorByStoreIDAndVendorID(ctx, nil, req.StoreID, &req.ID)
	if err != nil {
		vs.logger.Error("failed to get vendor by store ID and Vendor ID", zap.String("store_id", req.StoreID.String()), zap.String("vendorID", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get vendor by store ID and Vendor ID: %w", dto.ErrInternal)
	}
	if !found {
		vs.logger.Warn("vendor not found", zap.String("store_id", req.StoreID.String()), zap.String("vendorID", req.ID.String()))
		return nil, fmt.Errorf("vendor not found: %v", dto.ErrNotFound)
	}

	// validate email
	if req.Email != nil {
		if vendor.Email != *req.Email {
			_, found, err = vs.vendorRepo.GetVendorByStoreIDAndVendorEmail(ctx, nil, req.StoreID, *req.Email)
			if err != nil {
				vs.logger.Error("failed to get vendor by store ID and email", zap.String("store_id", req.StoreID.String()), zap.String("email", *req.Email), zap.Error(err))
				return nil, fmt.Errorf("failed to get vendor by store ID and email: %w", dto.ErrInternal)
			}
			if found {
				vs.logger.Warn("vendor email already exists", zap.String("email", *req.Email))
				return nil, fmt.Errorf("vendor already exists: %w", dto.ErrAlreadyExists)
			}
		}
		vendor.Email = *req.Email
	}

	// validate & normalize phone number
	phoneNumber, err := helper.NormalizePhoneNumber(req.PhoneNumber)
	if err != nil {
		vs.logger.Error("invalid phone number", zap.String("phone_number", req.PhoneNumber), zap.Error(err))
		return nil, fmt.Errorf("invalid phone number: %w", dto.ErrBadRequest)
	}
	if vendor.PhoneNumber != phoneNumber {
		_, found, err = vs.vendorRepo.GetVendorByStoreIDAndVendorPhoneNumber(ctx, nil, req.StoreID, phoneNumber)
		if err != nil {
			vs.logger.Error("failed to get vendor by store ID and phone number", zap.String("store_id", req.StoreID.String()), zap.String("phone_number", phoneNumber), zap.Error(err))
			return nil, fmt.Errorf("failed to get vendor by store ID and phone number: %w", dto.ErrInternal)
		}
		if found {
			vs.logger.Warn("vendor already exists", zap.String("phone_number", phoneNumber))
			return nil, fmt.Errorf("vendor already exists: %w", dto.ErrAlreadyExists)
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
	vendor.Name = req.Name

	err = vs.vendorRepo.UpdateVendor(ctx, nil, req.StoreID, vendor)
	if err != nil {
		vs.logger.Error("failed to update vendor", zap.String("id", vendor.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update vendor: %w", dto.ErrInternal)
	}

	return mapToVendorResponse(vendor), nil
}

func (vs *vendorService) DeleteVendorByStoreIDAndVendorID(ctx context.Context, storeID, vendorID *uuid.UUID) error {
	_, found, err := vs.vendorRepo.GetVendorByStoreIDAndVendorID(ctx, nil, storeID, vendorID)
	if err != nil {
		vs.logger.Error("failed to get vendor by store ID and Vendor ID", zap.String("store_id", storeID.String()), zap.String("vendorID", vendorID.String()), zap.Error(err))
		return fmt.Errorf("failed to get vendor by store ID and Vendor ID: %w", dto.ErrInternal)
	}
	if !found {
		vs.logger.Warn("vendor not found", zap.String("store_id", storeID.String()), zap.String("vendorID", vendorID.String()))
		return fmt.Errorf("vendor not found: %v", dto.ErrNotFound)
	}

	if err := vs.vendorRepo.DeleteVendorByStoreIDAndVendorID(ctx, nil, storeID, vendorID); err != nil {
		vs.logger.Error("failed to delete vendor by id", zap.String("store_id", storeID.String()), zap.String("vendorID", vendorID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete vendor by id: %w", dto.ErrInternal)
	}

	return nil
}
