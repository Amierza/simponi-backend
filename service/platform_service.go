// service/platform_service.go
package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	IPlatformService interface {
		// GetMyStore returns store + connected platforms milik user yang login.
		// Jika belum punya store, returns nil (bukan error).
		GetMyStore(ctx context.Context) (*dto.MyStoreResponse, error)

		// ConnectPlatform adalah mock OAuth:
		// - Jika user belum punya store → validasi store_name, buat Store + StoreUser + StorePlatform
		// - Jika sudah punya store → skip CreateStore, hanya tambah StorePlatform (cek duplikat)
		ConnectPlatform(ctx context.Context, req *dto.ConnectPlatformRequest) (*dto.MyStoreResponse, error)

		// DisconnectPlatform menghapus StorePlatform.
		// Jika setelah disconnect store tidak punya platform lagi → hapus store juga.
		DisconnectPlatform(ctx context.Context, storePlatformID *uuid.UUID) error
	}

	platformService struct {
		tx                repository.ITransaction
		storeRepo         repository.IStoreRepository
		storeUserRepo     repository.IStoreUserRepository
		platformRepo      repository.IPlatformRepository
		storePlatformRepo repository.IStorePlatformRepository
		logger            *zap.Logger
	}
)

func NewPlatformService(
	tx repository.ITransaction,
	storeRepo repository.IStoreRepository,
	storeUserRepo repository.IStoreUserRepository,
	platformRepo repository.IPlatformRepository,
	storePlatformRepo repository.IStorePlatformRepository,
	logger *zap.Logger,
) *platformService {
	return &platformService{
		tx:                tx,
		storeRepo:         storeRepo,
		storeUserRepo:     storeUserRepo,
		platformRepo:      platformRepo,
		storePlatformRepo: storePlatformRepo,
		logger:            logger,
	}
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func (ps *platformService) extractUserID(ctx context.Context) (uuid.UUID, error) {
	raw, ok := ctx.Value("user_id").(string)
	if !ok || raw == "" {
		return uuid.Nil, fmt.Errorf("missing user_id in context: %w", dto.ErrUnauthorized)
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user_id: %w", dto.ErrUnauthorized)
	}
	return id, nil
}

func buildMyStoreResponse(store *entity.Store) *dto.MyStoreResponse {
	res := &dto.MyStoreResponse{
		ID:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		ImageURL:    store.ImageURL,
		IsActive:    store.IsActive,
		Platforms:   []dto.ConnectedPlatformDetail{},
	}
	for _, sp := range store.StorePlatforms {
		if sp.Platform == nil {
			continue
		}
		res.Platforms = append(res.Platforms, dto.ConnectedPlatformDetail{
			StorePlatformID: sp.ID,
			PlatformID:      sp.Platform.ID,
			PlatformName:    sp.Platform.Name,
			ExternalName:    sp.ExternalName,
			ExternalShopID:  sp.ExternalShopID,
			IsConnected:     sp.IsConnected,
		})
	}
	return res
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ── GetMyStore ────────────────────────────────────────────────────────────────

func (ps *platformService) GetMyStore(ctx context.Context) (*dto.MyStoreResponse, error) {
	userID, err := ps.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	store, found, err := ps.storeRepo.GetStoreByUserID(ctx, nil, &userID)
	if err != nil {
		ps.logger.Error("failed to get store by user id", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get store: %w", dto.ErrInternal)
	}
	if !found {
		// Belum punya store — bukan error, FE handle kondisi ini
		return nil, nil
	}

	return buildMyStoreResponse(store), nil
}

// ── ConnectPlatform ───────────────────────────────────────────────────────────

func (ps *platformService) ConnectPlatform(ctx context.Context, req *dto.ConnectPlatformRequest) (*dto.MyStoreResponse, error) {
	userID, err := ps.extractUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Validasi platform exists di DB
	_, found, err := ps.platformRepo.GetPlatformByPlatformID(ctx, nil, req.PlatformID)
	if err != nil {
		ps.logger.Error("failed to get platform", zap.Error(err))
		return nil, fmt.Errorf("failed to validate platform: %w", dto.ErrInternal)
	}
	if !found {
		return nil, fmt.Errorf("platform not found: %w", dto.ErrNotFound)
	}

	// Cek apakah user sudah punya store
	existingStore, hasStore, err := ps.storeRepo.GetStoreByUserID(ctx, nil, &userID)
	if err != nil {
		ps.logger.Error("failed to get existing store", zap.Error(err))
		return nil, fmt.Errorf("failed to check existing store: %w", dto.ErrInternal)
	}

	if hasStore {
		// ── Path B: Store sudah ada → hanya tambah StorePlatform ─────────────
		//
		// Data store (nama, deskripsi, gambar) TIDAK diubah.
		// Hanya ExternalName & ExternalShopID yang dipakai dari request.

		targetStoreID := existingStore.ID

		// Cek duplikat platform
		_, alreadyConnected, err := ps.storePlatformRepo.GetStorePlatformByStoreIDAndPlatformID(
			ctx, nil, &targetStoreID, req.PlatformID,
		)
		if err != nil {
			ps.logger.Error("failed to check existing platform connection", zap.Error(err))
			return nil, fmt.Errorf("failed to check platform: %w", dto.ErrInternal)
		}
		if alreadyConnected {
			return nil, fmt.Errorf("platform already connected to this store: %w", dto.ErrAlreadyExists)
		}

		newStorePlatform := &entity.StorePlatform{
			ID:             uuid.New(),
			StoreID:        &targetStoreID,
			PlatformID:     req.PlatformID,
			ExternalShopID: req.ExternalShopID,
			ExternalName:   req.ExternalName,
			IsConnected:    true,
		}

		err = ps.tx.Run(ctx, func(tx *gorm.DB) error {
			return ps.storePlatformRepo.CreateStorePlatform(ctx, tx, newStorePlatform)
		})
		if err != nil {
			ps.logger.Error("failed to connect platform to existing store", zap.Error(err))
			return nil, fmt.Errorf("failed to connect platform: %w", dto.ErrInternal)
		}

		ps.logger.Info("second platform connected",
			zap.String("user_id", userID.String()),
			zap.String("store_id", targetStoreID.String()),
			zap.String("platform_id", req.PlatformID.String()),
		)

	} else {
		// ── Path A: Belum ada store → buat Store + StoreUser + StorePlatform ──

		// Validasi store_name wajib diisi saat platform pertama
		if req.StoreName == "" {
			return nil, fmt.Errorf("store_name is required for first platform connection: %w", dto.ErrBadRequest)
		}

		newStoreID := uuid.New()
		newStore := &entity.Store{
			ID:          newStoreID,
			Name:        req.StoreName,
			Description: derefString(req.StoreDescription),
			ImageURL:    derefString(req.StoreImageURL),
			IsActive:    true,
		}
		newStoreUser := &entity.StoreUser{
			ID:      uuid.New(),
			UserID:  &userID,
			StoreID: &newStoreID,
		}
		newStorePlatform := &entity.StorePlatform{
			ID:             uuid.New(),
			StoreID:        &newStoreID,
			PlatformID:     req.PlatformID,
			ExternalShopID: req.ExternalShopID,
			ExternalName:   req.ExternalName,
			IsConnected:    true,
		}

		err = ps.tx.Run(ctx, func(tx *gorm.DB) error {
			if err := ps.storeRepo.CreateStore(ctx, tx, newStore); err != nil {
				return err
			}
			if err := ps.storeUserRepo.CreateStoreUser(ctx, tx, newStoreUser); err != nil {
				return err
			}
			return ps.storePlatformRepo.CreateStorePlatform(ctx, tx, newStorePlatform)
		})
		if err != nil {
			ps.logger.Error("failed to create store and connect first platform", zap.Error(err))
			return nil, fmt.Errorf("failed to create store: %w", dto.ErrInternal)
		}

		ps.logger.Info("first platform connected — store created",
			zap.String("user_id", userID.String()),
			zap.String("store_id", newStoreID.String()),
			zap.String("platform_id", req.PlatformID.String()),
		)
	}

	// Reload store untuk response yang fresh
	updatedStore, _, err := ps.storeRepo.GetStoreByUserID(ctx, nil, &userID)
	if err != nil || updatedStore == nil {
		ps.logger.Error("failed to reload store after connect", zap.Error(err))
		return nil, fmt.Errorf("failed to reload store: %w", dto.ErrInternal)
	}

	return buildMyStoreResponse(updatedStore), nil
}

// ── DisconnectPlatform ────────────────────────────────────────────────────────

func (ps *platformService) DisconnectPlatform(ctx context.Context, storePlatformID *uuid.UUID) error {
	userID, err := ps.extractUserID(ctx)
	if err != nil {
		return err
	}

	// Validasi store_platform ada
	sp, found, err := ps.storePlatformRepo.GetStorePlatformByID(ctx, nil, storePlatformID)
	if err != nil {
		return fmt.Errorf("failed to get store platform: %w", dto.ErrInternal)
	}
	if !found {
		return fmt.Errorf("store platform not found: %w", dto.ErrNotFound)
	}

	// Validasi ownership
	store, found, err := ps.storeRepo.GetStoreByUserID(ctx, nil, &userID)
	if err != nil || !found {
		return fmt.Errorf("store not found for this user: %w", dto.ErrForbidden)
	}
	if store.ID != *sp.StoreID {
		return fmt.Errorf("access denied: %w", dto.ErrForbidden)
	}

	// Hitung sisa platform sebelum hapus
	count, err := ps.storePlatformRepo.CountStorePlatformsByStoreID(ctx, nil, sp.StoreID)
	if err != nil {
		return fmt.Errorf("failed to count platforms: %w", dto.ErrInternal)
	}

	err = ps.tx.Run(ctx, func(tx *gorm.DB) error {
		// Hapus store_platform
		if err := ps.storePlatformRepo.DeleteStorePlatformByID(ctx, tx, storePlatformID); err != nil {
			return err
		}

		// Jika ini platform terakhir → hapus store juga
		if count <= 1 {
			if err := ps.storeRepo.DeleteStoreByStoreID(ctx, tx, sp.StoreID); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		ps.logger.Error("failed to disconnect platform",
			zap.String("store_platform_id", storePlatformID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to disconnect platform: %w", dto.ErrInternal)
	}

	ps.logger.Info("platform disconnected",
		zap.String("user_id", userID.String()),
		zap.String("store_platform_id", storePlatformID.String()),
		zap.Bool("store_also_deleted", count <= 1),
	)
	return nil
}
