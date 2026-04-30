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
	IStoreUserService interface {
		CreateStoreUsers(ctx context.Context, req *dto.CreateStoreUsersRequest) error
		GetStoreUsers(ctx context.Context, req *response.PaginationRequest, storeID *uuid.UUID) (dto.StoreUserPaginationResponse, error)
		GetStoreUserByStoreIDAndUserID(ctx context.Context, storeID, userID *uuid.UUID) (*dto.UserResponse, error)
		DeleteStoreUserByStoreIDAndUserID(ctx context.Context, storeID, userID *uuid.UUID) error
	}

	storeUserService struct {
		storeUserRepo repository.IStoreUserRepository
		logger        *zap.Logger
		jwtService    jwt.IJWT
	}
)

func NewStoreUserService(storeUserRepo repository.IStoreUserRepository, logger *zap.Logger, jwtService jwt.IJWT) *storeUserService {
	return &storeUserService{
		storeUserRepo: storeUserRepo,
		logger:        logger,
		jwtService:    jwtService,
	}
}

func (sus *storeUserService) CreateStoreUsers(ctx context.Context, req *dto.CreateStoreUsersRequest) error {
	// check if user already exist in this store
	for _, userID := range req.UserIDs {
		_, found, err := sus.storeUserRepo.GetStoreUserByStoreIDAndUserID(ctx, nil, req.StoreID, userID)
		if err != nil {
			sus.logger.Error("failed to get store user by store id and user id", zap.String("store_id", req.StoreID.String()), zap.String("user_id", userID.String()), zap.Error(err))
			return fmt.Errorf("failed to get store user by store id and user id: %w", dto.ErrInternal)
		}
		if found {
			sus.logger.Warn("user already exists in store", zap.String("user_id", userID.String()))
			return fmt.Errorf("user already exists in store: %w", dto.ErrAlreadyExists)
		}

		newID := uuid.New()
		newStoreUser := &entity.StoreUser{
			ID:      newID,
			StoreID: req.StoreID,
			UserID:  userID,
		}

		err = sus.storeUserRepo.CreateStoreUser(ctx, nil, newStoreUser)
		if err != nil {
			sus.logger.Error("failed to create storeUser", zap.Error(err))
			return fmt.Errorf("failed to create storeUser: %w", dto.ErrInternal)
		}
	}

	sus.logger.Info("success to create store users")

	return nil
}

func (sus *storeUserService) GetStoreUsers(ctx context.Context, req *response.PaginationRequest, storeID *uuid.UUID) (dto.StoreUserPaginationResponse, error) {
	datas, err := sus.storeUserRepo.GetStoreUsersByStoreID(ctx, nil, req, storeID)
	if err != nil {
		sus.logger.Error("failed to get store users by store id", zap.Error(err))
		return dto.StoreUserPaginationResponse{}, fmt.Errorf("failed to get store users by store id: %w", dto.ErrInternal)
	}

	sus.logger.Info("success to get store users by store id", zap.String("store_id", storeID.String()), zap.Int64("count", datas.Count))

	var storeUsers []*dto.CustomUserResponse
	for _, storeUser := range datas.StoreUsers {
		storeUsers = append(storeUsers, &dto.CustomUserResponse{
			ID:    storeUser.User.ID,
			Name:  storeUser.User.Name,
			Email: storeUser.User.Email,
		})
	}

	return dto.StoreUserPaginationResponse{
		Data:               storeUsers,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (sus *storeUserService) GetStoreUserByStoreIDAndUserID(ctx context.Context, storeID, userID *uuid.UUID) (*dto.UserResponse, error) {
	storeUser, found, err := sus.storeUserRepo.GetStoreUserByStoreIDAndUserID(ctx, nil, storeID, userID)
	if err != nil {
		sus.logger.Error("failed to get store user by store id and user id", zap.String("store_id", storeID.String()), zap.String("user_id", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get store user: %w", dto.ErrInternal)
	}
	if !found {
		sus.logger.Warn("store user not found", zap.String("store_id", storeID.String()), zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("store user not found: %v", dto.ErrNotFound)
	}

	sus.logger.Info("success to get store user by store_id", zap.String("store_id", storeID.String()))

	return &dto.UserResponse{
		ID:       storeUser.User.ID,
		Email:    storeUser.User.Email,
		Name:     storeUser.User.Name,
		ImageURL: storeUser.User.ImageURL,
		Status:   storeUser.User.Status,
		Role: dto.RoleResponse{
			ID:   storeUser.User.Role.ID,
			Name: storeUser.User.Role.Name,
		},
	}, nil
}

func (sus *storeUserService) DeleteStoreUserByStoreIDAndUserID(ctx context.Context, storeID, userID *uuid.UUID) error {
	_, found, err := sus.storeUserRepo.GetStoreUserByStoreIDAndUserID(ctx, nil, storeID, userID)
	if err != nil {
		sus.logger.Error("failed to get store user by store id and user id", zap.String("store_id", storeID.String()), zap.String("user_id", userID.String()), zap.Error(err))
		return fmt.Errorf("failed to get store user: %w", dto.ErrInternal)
	}
	if !found {
		sus.logger.Warn("store user not found", zap.String("store_id", storeID.String()), zap.String("user_id", userID.String()))
		return fmt.Errorf("store user not found: %v", dto.ErrNotFound)
	}

	if err := sus.storeUserRepo.DeleteStoreUserByStoreIDAndUserID(ctx, nil, storeID, userID); err != nil {
		sus.logger.Error("failed to delete store user by store id and user id", zap.String("store_id", storeID.String()), zap.String("user_id", userID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete store user by store id and user id: %w", dto.ErrInternal)
	}

	return nil
}
