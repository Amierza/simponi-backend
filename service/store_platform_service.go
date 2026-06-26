// service/platform_service.go
package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	IStorePlatformService interface {
		ConnectPlatform(ctx context.Context, storeID, platformID *uuid.UUID) (dto.ConnectPlatformResponse, error)
		DisconnectPlatform(ctx context.Context, storeID, platformID *uuid.UUID) error
		SyncProducts(ctx context.Context, storeID, platformID *uuid.UUID) error
	}

	storePlatformService struct {
		tx                  repository.ITransaction
		storeRepo           repository.IStoreRepository
		platformRepo        repository.IPlatformRepository
		storePlatformRepo   repository.IStorePlatformRepository
		storeCredentialRepo repository.IStoreCredentialRepository
		logger              *zap.Logger
		shopeeService       *shopeeService
	}
)

func NewStorePlatformService(
	tx repository.ITransaction,
	storePlatformRepo repository.IStorePlatformRepository,
	storeRepo repository.IStoreRepository,
	platformRepo repository.IPlatformRepository,
	storeCredentialRepo repository.IStoreCredentialRepository,
	logger *zap.Logger,
	shopeeService *shopeeService,
) *storePlatformService {
	return &storePlatformService{
		tx:                  tx,
		storeRepo:           storeRepo,
		platformRepo:        platformRepo,
		storePlatformRepo:   storePlatformRepo,
		storeCredentialRepo: storeCredentialRepo,
		logger:              logger,
		shopeeService:       shopeeService,
	}
}

func (sps *storePlatformService) ConnectPlatform(ctx context.Context, storeID, platformID *uuid.UUID) (dto.ConnectPlatformResponse, error) {
	platform, found, err := sps.platformRepo.GetPlatformByPlatformID(ctx, nil, platformID)
	if err != nil {
		sps.logger.Error("failed to get platform by ID", zap.String("platformID", platformID.String()), zap.Error(err))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("failed to get platform ID: %w", dto.ErrInternal)
	}
	if !found {
		sps.logger.Warn("platform not found", zap.String("platformID", platformID.String()))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("platform not found: %v", dto.ErrNotFound)
	}

	_, found, err = sps.storeRepo.GetStoreByStoreID(ctx, nil, storeID)
	if err != nil {
		sps.logger.Error("failed to get store by ID", zap.String("storeID", storeID.String()), zap.Error(err))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("failed to get store ID: %w", dto.ErrInternal)
	}
	if !found {
		sps.logger.Warn("store not found", zap.String("storeID", storeID.String()))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("store not found: %v", dto.ErrNotFound)
	}

	if platformConnectMode() == "dummy" {
		return sps.connectPlatformDummy(ctx, storeID, platform)
	}

	_, found, err = sps.storePlatformRepo.GetStorePlatformByStoreIDAndPlatformID(ctx, nil, storeID, platformID)
	if err != nil {
		sps.logger.Error("failed to get store platform by IDs", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()), zap.Error(err))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("failed to get store platform: %w", dto.ErrInternal)
	}
	if found {
		sps.logger.Error("store platform already exists", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()), zap.Error(err))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("store platform already exists: %w", dto.ErrAlreadyExists)
	}

	switch strings.ToLower(platform.Name) {
	case "shopee":
		authURL, err := sps.shopeeService.GenerateAuthURL(ctx, storeID, platformID)
		if err != nil {
			sps.logger.Error("failed generate shopee auth url", zap.Error(err), zap.String("store_id", storeID.String()), zap.String("platform_id", platformID.String()))
			return dto.ConnectPlatformResponse{}, fmt.Errorf("failed to connect platform: %w", dto.ErrInternal)
		}

		sps.logger.Info("success generate shopee auth url", zap.String("store_id", storeID.String()), zap.String("platform_id", platformID.String()))

		return dto.ConnectPlatformResponse{
			AuthURL: authURL,
		}, nil

	default:
		sps.logger.Error("platform not supported", zap.String("platformID", platformID.String()))
		return dto.ConnectPlatformResponse{}, fmt.Errorf("platform not supported: %w", dto.ErrBadRequest)
	}
}

func platformConnectMode() string {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("PLATFORM_CONNECT_MODE")))
	if mode == "" {
		return "oauth"
	}

	return mode
}

func dummyPlatformName(platformName string) string {
	name := strings.ToLower(strings.TrimSpace(platformName))
	name = strings.ReplaceAll(name, " ", "-")
	if name == "" {
		return "platform"
	}

	return name
}

func (sps *storePlatformService) connectPlatformDummy(ctx context.Context, storeID *uuid.UUID, platform *entity.Platform) (dto.ConnectPlatformResponse, error) {
	platformName := dummyPlatformName(platform.Name)
	storeIDString := storeID.String()
	shortStoreID := storeIDString
	if len(shortStoreID) > 8 {
		shortStoreID = shortStoreID[:8]
	}

	platformID := platform.ID
	externalShopID := fmt.Sprintf("dummy-%s-%s", platformName, shortStoreID)
	externalName := fmt.Sprintf("Dummy %s Store", platformName)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	var storePlatform *entity.StorePlatform
	err := sps.tx.Run(ctx, func(tx *gorm.DB) error {
		upsertedStorePlatform, err := sps.storePlatformRepo.UpsertStorePlatformByStoreIDAndPlatformID(ctx, tx, &entity.StorePlatform{
			ID:             uuid.New(),
			StoreID:        storeID,
			PlatformID:     &platformID,
			ExternalShopID: externalShopID,
			ExternalName:   externalName,
			IsConnected:    true,
		})
		if err != nil {
			sps.logger.Error("failed upsert dummy store platform", zap.String("store_id", storeIDString), zap.String("platform_id", platformID.String()), zap.Error(err))
			return fmt.Errorf("failed to connect platform: %w", dto.ErrInternal)
		}

		storePlatform = upsertedStorePlatform

		_, err = sps.storeCredentialRepo.UpsertStoreCredentialByStorePlatformID(ctx, tx, &entity.StoreCredential{
			ID:              uuid.New(),
			StorePlatformID: storePlatform.ID,
			AccessToken:     fmt.Sprintf("dummy_access_token_%s_%s", platformName, storeIDString),
			RefreshToken:    fmt.Sprintf("dummy_refresh_token_%s_%s", platformName, storeIDString),
			ExpiresAt:       &expiresAt,
		})
		if err != nil {
			sps.logger.Error("failed upsert dummy store credential", zap.String("store_platform_id", storePlatform.ID.String()), zap.Error(err))
			return fmt.Errorf("failed to connect platform: %w", dto.ErrInternal)
		}

		return nil
	})
	if err != nil {
		return dto.ConnectPlatformResponse{}, err
	}

	sps.logger.Warn("dummy platform connect mode active", zap.String("store_id", storeIDString), zap.String("platform_id", platformID.String()), zap.String("platform_name", platform.Name))

	return dto.ConnectPlatformResponse{
		StorePlatformID: &storePlatform.ID,
		PlatformID:      &platformID,
		PlatformName:    platform.Name,
		ExternalShopID:  storePlatform.ExternalShopID,
		ExternalName:    storePlatform.ExternalName,
		IsConnected:     &storePlatform.IsConnected,
	}, nil
}

func (sps *storePlatformService) DisconnectPlatform(ctx context.Context, storeID, platformID *uuid.UUID) error {
	storePlatform, found, err := sps.storePlatformRepo.GetStorePlatformByStoreIDAndPlatformID(ctx, nil, storeID, platformID)
	if err != nil {
		sps.logger.Error("failed to get store platform by IDs", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()), zap.Error(err))
		return fmt.Errorf("failed to get store platform: %w", dto.ErrInternal)
	}
	if !found {
		sps.logger.Warn("store platform not found", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()))
		return fmt.Errorf("store platform not found: %v", dto.ErrNotFound)
	}

	err = sps.tx.Run(ctx, func(tx *gorm.DB) error {
		// delete credential if exists
		if storePlatform.Credential != nil {
			err := sps.storeCredentialRepo.DeleteStoreCredentialByID(ctx, tx, &storePlatform.Credential.ID)
			if err != nil {
				sps.logger.Error("failed to delete store credential", zap.String("credentialID", storePlatform.Credential.ID.String()), zap.Error(err))
				return fmt.Errorf("failed to disconnect platform: %w", dto.ErrInternal)
			}
			sps.logger.Info("store credential deleted successfully", zap.String("credentialID", storePlatform.Credential.ID.String()))
		} else {
			sps.logger.Info("no credential to delete for this store platform", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()))
		}

		err = sps.storePlatformRepo.DeleteStorePlatformByID(ctx, tx, &storePlatform.ID)
		if err != nil {
			sps.logger.Error("failed to delete store platform", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()), zap.Error(err))
			return fmt.Errorf("failed to disconnect platform: %w", dto.ErrInternal)
		}

		return nil
	})

	sps.logger.Info("store platform disconnected successfully", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()))
	return nil
}

func (sps *storePlatformService) SyncProducts(ctx context.Context, storeID, platformID *uuid.UUID) error {
	storePlatform, found, err := sps.storePlatformRepo.GetStorePlatformByStoreIDAndPlatformID(ctx, nil, storeID, platformID)
	if err != nil {
		sps.logger.Error("failed to get store platform by IDs", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()), zap.Error(err))
		return fmt.Errorf("failed to get store platform: %w", dto.ErrInternal)
	}
	if !found {
		sps.logger.Warn("store platform not found", zap.String("storeID", storeID.String()), zap.String("platformID", platformID.String()))
		return fmt.Errorf("store platform not found: %v", dto.ErrNotFound)
	}

	switch strings.ToLower(storePlatform.Platform.Name) {
	case "shopee":
		err := sps.shopeeService.SyncProducts(ctx, storePlatform)
		if err != nil {
			sps.logger.Error("failed to sync products from shopee", zap.String("storePlatformID", storePlatform.ID.String()), zap.Error(err))
			return fmt.Errorf("failed to sync products: %w", dto.ErrInternal)
		}
		sps.logger.Info("products synced successfully from shopee", zap.String("storePlatformID", storePlatform.ID.String()))
		return nil

	default:
		sps.logger.Error("platform not supported for syncing products", zap.String("platformName", storePlatform.Platform.Name))
		return fmt.Errorf("platform not supported for syncing products: %w", dto.ErrBadRequest)
	}
}
