package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	IDashboardRepository interface {
		GetDashboardData(ctx context.Context, tx *gorm.DB) (dto.DashboardSuperadminResponse, error)
	}

	dashboardRepository struct {
		db *gorm.DB
	}
)

func NewDashboardRepository(db *gorm.DB) *dashboardRepository {
	return &dashboardRepository{
		db: db,
	}
}

func (dr *dashboardRepository) GetDashboardData(ctx context.Context, tx *gorm.DB) (dto.DashboardSuperadminResponse, error) {
	if tx == nil {
		tx = dr.db
	}

	var response dto.DashboardSuperadminResponse

	// Total Users
	if err := tx.WithContext(ctx).
		Model(&entity.User{}).
		Count(&response.TotalUsers).Error; err != nil {
		return dto.DashboardSuperadminResponse{}, err
	}

	// Total Stores
	if err := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Count(&response.TotalStores).Error; err != nil {
		return dto.DashboardSuperadminResponse{}, err
	}

	// Total Connected TikTok Shop
	if err := tx.WithContext(ctx).
		Model(&entity.StorePlatform{}).
		Joins("JOIN platforms p ON p.id = store_platforms.platform_id").
		Where("store_platforms.is_connected = ?", true).
		Where("LOWER(p.name) = ?", "tiktok shop").
		Count(&response.TotalConnectedTikTokShop).Error; err != nil {
		return dto.DashboardSuperadminResponse{}, err
	}

	// Total Connected Shopee
	if err := tx.WithContext(ctx).
		Model(&entity.StorePlatform{}).
		Joins("JOIN platforms p ON p.id = store_platforms.platform_id").
		Where("store_platforms.is_connected = ?", true).
		Where("LOWER(p.name) = ?", "shopee").
		Count(&response.TotalConnectedShopee).Error; err != nil {
		return dto.DashboardSuperadminResponse{}, err
	}

	// New Stores Last 7 Days
	if err := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Where("created_at >= NOW() - INTERVAL '7 days'").
		Count(&response.NewStoresLast7Days).Error; err != nil {
		return dto.DashboardSuperadminResponse{}, err
	}

	return response, nil
}
