package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/constants"
	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	IStoreService interface {
		CreateStore(ctx context.Context, req *dto.CreateStoreRequest) (*dto.StoreResponse, error)
		GetStores(ctx context.Context, req *response.PaginationRequest) (dto.StorePaginationResponse, error)
		GetStoreByStoreID(ctx context.Context, storeID *uuid.UUID) (*dto.StoreResponse, error)
		UpdateStoreByStoreID(ctx context.Context, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error)
		DeleteStoreByStoreID(ctx context.Context, storeID *uuid.UUID) error
	}

	storeService struct {
		tx                repository.ITransaction
		storeRepo         repository.IStoreRepository
		storeUserRepo     repository.IStoreUserRepository
		platformRepo      repository.IPlatformRepository
		storePlatformRepo repository.IStorePlatformRepository
		logger            *zap.Logger
		jwtService        jwt.IJWT
	}
)

func NewStoreService(tx repository.ITransaction, storeRepo repository.IStoreRepository, storeUserRepo repository.IStoreUserRepository, platformRepo repository.IPlatformRepository, storePlatformRepo repository.IStorePlatformRepository, logger *zap.Logger, jwtService jwt.IJWT) *storeService {
	return &storeService{
		tx:                tx,
		storeRepo:         storeRepo,
		storeUserRepo:     storeUserRepo,
		platformRepo:      platformRepo,
		storePlatformRepo: storePlatformRepo,
		logger:            logger,
		jwtService:        jwtService,
	}
}

func mapToStoreResponse(s *entity.Store) *dto.StoreResponse {
	res := &dto.StoreResponse{
		ID:          s.ID,
		Name:        s.Name,
		ImageURL:    s.ImageURL,
		Description: s.Description,
		IsActive:    s.IsActive,
	}

	for _, p := range s.StorePlatforms {
		res.Platforms = append(res.Platforms, dto.PlatformResponse{
			ID:   p.Platform.ID,
			Name: p.Platform.Name,
		})
	}

	return res
}

func (ss *storeService) CreateStore(ctx context.Context, req *dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	userIDString := ctx.Value("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		ss.logger.Error("failed to parse user_id", zap.String("user_id", userIDString), zap.Error(err))
		return nil, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
	}

	_, found, err := ss.platformRepo.GetPlatformByPlatformID(ctx, nil, req.PlatformID)
	if err != nil {
		ss.logger.Error("failed to get platform by ID", zap.String("platformID", req.PlatformID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get platform ID: %w", dto.ErrInternal)
	}
	if !found {
		ss.logger.Warn("platform not found", zap.String("platformID", req.PlatformID.String()))
		return nil, fmt.Errorf("platform not found: %v", dto.ErrNotFound)
	}

	newStoreID := uuid.New()
	newStore := &entity.Store{
		ID:          newStoreID,
		Name:        req.Name,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		IsActive:    true,
	}

	newStoreUserID := uuid.New()
	newStoreUser := &entity.StoreUser{
		ID:      newStoreUserID,
		UserID:  &userID,
		StoreID: &newStoreID,
	}

	newStorePlatformID := uuid.New()
	newStorePlatform := &entity.StorePlatform{
		ID:         newStorePlatformID,
		StoreID:    &newStoreID,
		PlatformID: req.PlatformID,
	}

	err = ss.tx.Run(ctx, func(tx *gorm.DB) error {
		err = ss.storeRepo.CreateStore(ctx, tx, newStore)
		if err != nil {
			ss.logger.Error("failed to create store", zap.Error(err))
			return fmt.Errorf("failed to create store: %w", dto.ErrInternal)
		}

		err = ss.storeUserRepo.CreateStoreUser(ctx, tx, newStoreUser)
		if err != nil {
			ss.logger.Error("failed to create store user", zap.Error(err))
			return fmt.Errorf("failed to create store user: %w", dto.ErrInternal)
		}

		err = ss.storePlatformRepo.CreateStorePlatform(ctx, tx, newStorePlatform)
		if err != nil {
			ss.logger.Error("failed to create store platform", zap.Error(err))
			return fmt.Errorf("failed to create store platform: %w", dto.ErrInternal)
		}
		return nil
	})

	store, found, err := ss.storeRepo.GetStoreByStoreID(ctx, nil, &newStoreID)
	if err != nil {
		ss.logger.Error("failed to get store by ID", zap.String("storeID", newStoreID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}

	ss.logger.Info("success to create store", zap.String("id", newStore.ID.String()))

	return mapToStoreResponse(store), nil
}

func (ss *storeService) GetStores(ctx context.Context, req *response.PaginationRequest) (dto.StorePaginationResponse, error) {
	userIDString := ctx.Value("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		ss.logger.Error("failed to parse user_id", zap.String("user_id", userIDString), zap.Error(err))
		return dto.StorePaginationResponse{}, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
	}

	roleIDString := ctx.Value("role_id").(string)
	var datas dto.StorePaginationRepositoryResponse

	if roleIDString == constants.SUPER_ADMIN_ROLE_ID {
		datas, err = ss.storeRepo.GetStores(ctx, nil, req)
		if err != nil {
			ss.logger.Error("failed to get stores", zap.Error(err))
			return dto.StorePaginationResponse{}, fmt.Errorf("failed to get stores: %w", dto.ErrInternal)
		}
	} else {
		datas, err = ss.storeRepo.GetStoresByUserID(ctx, nil, req, &userID)
		if err != nil {
			ss.logger.Error("failed to get stores by user id", zap.String("user_id", userIDString), zap.Error(err))
			return dto.StorePaginationResponse{}, fmt.Errorf("failed to get stores: %w", dto.ErrInternal)
		}
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

func (ss *storeService) GetStoreByStoreID(ctx context.Context, storeID *uuid.UUID) (*dto.StoreResponse, error) {
	store, found, err := ss.storeRepo.GetStoreByStoreID(ctx, nil, storeID)
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

func (ss *storeService) UpdateStoreByStoreID(ctx context.Context, req *dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	store, found, err := ss.storeRepo.GetStoreByStoreID(ctx, nil, &req.ID)
	if err != nil {
		ss.logger.Error("failed to get store by ID", zap.String("storeID", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}
	if !found {
		ss.logger.Warn("store not found", zap.String("storeID", req.ID.String()))
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

	err = ss.storeRepo.UpdateStoreByStoreID(ctx, nil, store)
	if err != nil {
		ss.logger.Error("failed to update store", zap.String("id", req.ID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to update store: %w", dto.ErrInternal)
	}

	return mapToStoreResponse(store), nil
}

func (ss *storeService) DeleteStoreByStoreID(ctx context.Context, storeID *uuid.UUID) error {
	_, found, err := ss.storeRepo.GetStoreByStoreID(ctx, nil, storeID)
	if err != nil {
		ss.logger.Error("failed to get store by ID", zap.String("storeID", storeID.String()), zap.Error(err))
		return fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}
	if !found {
		ss.logger.Warn("store not found", zap.String("storeID", storeID.String()))
		return fmt.Errorf("store not found: %v", dto.ErrNotFound)
	}

	if err := ss.storeRepo.DeleteStoreByStoreID(ctx, nil, storeID); err != nil {
		ss.logger.Error("failed to delete store by id", zap.String("storeID", storeID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete store by id: %w", dto.ErrInternal)
	}

	return nil
}
