package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IStoreService interface {
		CreateStore(ctx context.Context, req *dto.CreateStoreRequest) (*dto.StoreResponse, error)
		GetStores(ctx context.Context, req *response.PaginationRequest) (dto.StorePaginationResponse, error)
		GetStoreByID(ctx context.Context, storeID *uuid.UUID) (*dto.StoreResponse, error)
		UpdateStore(ctx context.Context, storeID *uuid.UUID, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error)
		DeleteStoreByID(ctx context.Context, storeID *uuid.UUID) error
	}

	storeService struct {
		storeRepo  repository.IStoreRepository
		logger     *zap.Logger
		jwtService jwt.IJWT
	}
)

func NewStoreService(storeRepo repository.IStoreRepository, logger *zap.Logger, jwtService jwt.IJWT) *storeService {
	return &storeService{
		storeRepo:  storeRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}

func mapToStoreResponse(v *entity.Store) *dto.StoreResponse {
	return &dto.StoreResponse{
		ID:          v.ID,
		Name:        v.Name,
		ImageURL:    v.ImageURL,
		Description: v.Description,
		IsActive:    v.IsActive,
	}
}

func (ss *storeService) CreateStore(ctx context.Context, req *dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	userIDString := ctx.Value("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		ss.logger.Error("failed to parse user_id", zap.String("user_id", userIDString), zap.Error(err))
		return nil, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
	}

	newID := uuid.New()
	newStore := &entity.Store{
		ID:          newID,
		Name:        req.Name,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		IsActive:    true,
		UserID:      &userID,
	}

	err = ss.storeRepo.CreateStore(ctx, nil, newStore)
	if err != nil {
		ss.logger.Error("failed to create store", zap.Error(err))
		return nil, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
	}

	ss.logger.Info("success to create store", zap.String("id", newStore.ID.String()))

	return mapToStoreResponse(newStore), nil
}

func (ss *storeService) GetStores(ctx context.Context, req *response.PaginationRequest) (dto.StorePaginationResponse, error) {
	datas, err := ss.storeRepo.GetStores(ctx, nil, req)
	if err != nil {
		ss.logger.Error("failed to get stores", zap.Error(err))
		return dto.StorePaginationResponse{}, fmt.Errorf("failed to get stores: %w", dto.ErrInternal)
	}

	ss.logger.Info("success to get stores", zap.Int64("count", datas.Count))

	var stores []*dto.StoreResponse
	for _, store := range datas.Stores {
		stores = append(stores, mapToStoreResponse(store))
	}

	return dto.StorePaginationResponse{
		Data:               stores,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (ss *storeService) GetStoreByID(ctx context.Context, storeID *uuid.UUID) (*dto.StoreResponse, error) {
	store, found, err := ss.storeRepo.GetStoreByID(ctx, nil, storeID)
	if err != nil {
		ss.logger.Error("failed to get store by ID", zap.String("storeID", storeID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}
	if !found {
		ss.logger.Warn("store not found", zap.String("storeID", storeID.String()))
		return nil, fmt.Errorf("store not found: %v", dto.ErrNotFound)
	}

	ss.logger.Info("success to get store by id", zap.String("id", storeID.String()))

	return mapToStoreResponse(store), nil
}

func (ss *storeService) UpdateStore(ctx context.Context, storeID *uuid.UUID, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	store, found, err := ss.storeRepo.GetStoreByID(ctx, nil, storeID)
	if err != nil {
		ss.logger.Error("failed to get store by ID", zap.String("storeID", storeID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}
	if !found {
		ss.logger.Warn("store not found", zap.String("storeID", storeID.String()))
		return nil, fmt.Errorf("store not found: %v", dto.ErrNotFound)
	}

	if req.ImageURL != nil {
		store.ImageURL = *req.ImageURL
	}
	if req.Description != nil {
		store.Description = *req.Description
	}
	if req.IsActive != nil {
		store.IsActive = *req.IsActive
	}
	store.Name = req.Name

	err = ss.storeRepo.UpdateStore(ctx, nil, store)
	if err != nil {
		ss.logger.Error("failed to update store", zap.String("id", storeID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update store: %w", dto.ErrInternal)
	}

	return mapToStoreResponse(store), nil
}

func (ss *storeService) DeleteStoreByID(ctx context.Context, storeID *uuid.UUID) error {
	_, found, err := ss.storeRepo.GetStoreByID(ctx, nil, storeID)
	if err != nil {
		ss.logger.Error("failed to get store by ID", zap.String("storeID", storeID.String()), zap.Error(err))
		return fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}
	if !found {
		ss.logger.Warn("store not found", zap.String("storeID", storeID.String()))
		return fmt.Errorf("store not found: %v", dto.ErrNotFound)
	}

	if err := ss.storeRepo.DeleteStoreByID(ctx, nil, storeID); err != nil {
		ss.logger.Error("failed to delete store by id", zap.String("storeID", storeID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete store by id: %w", dto.ErrInternal)
	}

	return nil
}
