package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	IOAuthService interface {
		HandleShopeeCallback(ctx context.Context, code, shopID, state string) error
	}

	oAuthService struct {
		tx                  repository.ITransaction
		storePlatformRepo   repository.IStorePlatformRepository
		storeCredentialRepo repository.IStoreCredentialRepository
		redis               *redis.Client
		logger              *zap.Logger
		shopeeService       IShopeeService
	}
)

func NewOAuthService(tx repository.ITransaction, storePlatformRepo repository.IStorePlatformRepository, storeCredentialRepo repository.IStoreCredentialRepository, redis *redis.Client, logger *zap.Logger, shopeeService IShopeeService) *oAuthService {
	return &oAuthService{
		tx:                  tx,
		storePlatformRepo:   storePlatformRepo,
		storeCredentialRepo: storeCredentialRepo,
		redis:               redis,
		logger:              logger,
		shopeeService:       shopeeService,
	}
}

func (oas *oAuthService) HandleShopeeCallback(ctx context.Context, code, shopID, state string) error {
	if code == "" || shopID == "" || state == "" {
		oas.logger.Error("missing shopee callback query params", zap.String("code", code), zap.String("shop_id", shopID), zap.String("state", state))
		return fmt.Errorf("invalid oauth callback params: %w", dto.ErrBadRequest)
	}

	// get oauth state from redis
	redisKey := fmt.Sprintf("oauth:shopee:state:%s", state)
	stateData, err := oas.redis.Get(ctx, redisKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			oas.logger.Warn("oauth state not found", zap.String("state", state))
			return fmt.Errorf("oauth state expired or invalid: %w", dto.ErrUnauthorized)
		}

		oas.logger.Error("failed get oauth state from redis", zap.String("state", state), zap.Error(err))
		return fmt.Errorf("failed get oauth state: %w", dto.ErrInternal)
	}

	var payload dto.ShopeeOAuthState
	if err := json.Unmarshal([]byte(stateData), &payload); err != nil {
		oas.logger.Error("failed unmarshal oauth state", zap.String("state", state), zap.Error(err))
		return fmt.Errorf("failed parse oauth state: %w", dto.ErrInternal)
	}

	storeID, err := uuid.Parse(payload.StoreID)
	if err != nil {
		oas.logger.Error("failed parse store id", zap.String("store_id", payload.StoreID), zap.Error(err))
		return fmt.Errorf("failed parse store id: %w", dto.ErrInternal)
	}

	platformID, err := uuid.Parse(payload.PlatformID)
	if err != nil {
		oas.logger.Error("failed parse platform id", zap.String("platform_id", payload.PlatformID), zap.Error(err))
		return fmt.Errorf("failed parse platform id: %w", dto.ErrInternal)
	}

	// exchange code -> access token
	tokenResp, err := oas.shopeeService.GetAccessToken(ctx, code, shopID)
	if err != nil {
		oas.logger.Error("failed get shopee access token", zap.Error(err), zap.String("shop_id", shopID))
		return fmt.Errorf("failed get access token: %w", dto.ErrInternal)
	}

	shopIDInt, err := strconv.Atoi(shopID)
	if err != nil {
		oas.logger.Error("failed parse shop id", zap.String("shop_id", shopID), zap.Error(err))
		return fmt.Errorf("failed parse shop id: %w", dto.ErrInternal)
	}

	expiredAt := time.Now().Add(time.Duration(tokenResp.ExpireIn) * time.Second)
	err = oas.tx.Run(ctx, func(tx *gorm.DB) error {
		// store platform
		storePlatform := &entity.StorePlatform{
			ID:             uuid.New(),
			StoreID:        &storeID,
			PlatformID:     &platformID,
			ExternalShopID: strconv.Itoa(shopIDInt),
			IsConnected:    true,
		}
		err = oas.storePlatformRepo.CreateStorePlatform(ctx, tx, storePlatform)
		if err != nil {
			oas.logger.Error("failed create store platform", zap.Error(err))
			return fmt.Errorf("failed create store platform: %w", dto.ErrInternal)
		}

		// store credential
		credential := &entity.StoreCredential{
			ID:              uuid.New(),
			StorePlatformID: storePlatform.ID,
			AccessToken:     tokenResp.AccessToken,
			RefreshToken:    tokenResp.RefreshToken,
			ExpiresAt:       &expiredAt,
		}
		err = oas.storeCredentialRepo.CreateStoreCredential(ctx, tx, credential)
		if err != nil {
			oas.logger.Error("failed create store credential", zap.Error(err))
			return fmt.Errorf("failed create store credential: %w", dto.ErrInternal)
		}

		return nil
	})
	if err != nil {
		oas.logger.Error("failed to connect shopee platform", zap.Error(err))
		return fmt.Errorf("failed to connect shopee platform: %w", dto.ErrInternal)
	}

	// delete oauth state
	if err := oas.redis.Del(ctx, redisKey).Err(); err != nil {
		oas.logger.Warn("failed delete oauth state", zap.String("redis_key", redisKey), zap.Error(err))
	}

	oas.logger.Info("success connect shopee platform", zap.String("store_id", storeID.String()), zap.String("platform_id", platformID.String()), zap.String("external_shop_id", shopID))
	return nil
}
