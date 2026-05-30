package service

import (
	"context"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
    IDashboardService interface {
        GetDashboard(ctx context.Context, storeID *uuid.UUID) (dto.DashboardResponse, error)
    }

    dashboardService struct {
        dashboardRepo repository.IDashboardRepository
        logger        *zap.Logger
    }
)

func NewDashboardService(dashboardRepo repository.IDashboardRepository, logger *zap.Logger) IDashboardService {
    return &dashboardService{dashboardRepo: dashboardRepo, logger: logger}
}

func (ds *dashboardService) GetDashboard(ctx context.Context, storeID *uuid.UUID) (dto.DashboardResponse, error) {
    var res dto.DashboardResponse

    now := time.Now()
    from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
    to := now

    summary, err := ds.dashboardRepo.GetSummary(ctx, nil, storeID, from, to)
    if err != nil {
        ds.logger.Error("failed to get dashboard summary", zap.Error(err))
        return res, err
    }

    trend, err := ds.dashboardRepo.GetTrend(ctx, nil, storeID, 12)
    if err != nil {
        ds.logger.Error("failed to get dashboard trend", zap.Error(err))
        return res, err
    }

    recentOrders, err := ds.dashboardRepo.GetRecentOrders(ctx, nil, storeID, 5)
    if err != nil {
        ds.logger.Error("failed to get recent orders", zap.Error(err))
        return res, err
    }

    lowStock, err := ds.dashboardRepo.GetLowStock(ctx, nil, storeID, 10, 10)
    if err != nil {
        ds.logger.Error("failed to get low stock", zap.Error(err))
        return res, err
    }

    topProducts, err := ds.dashboardRepo.GetTopProducts(ctx, nil, storeID, 5)
    if err != nil {
        ds.logger.Error("failed to get top products", zap.Error(err))
        return res, err
    }

    activity, err := ds.dashboardRepo.GetActivity(ctx, nil, storeID, 10)
    if err != nil {
        ds.logger.Error("failed to get activity", zap.Error(err))
        return res, err
    }

    res = dto.DashboardResponse{
        Summary: dto.DashboardSummaryResponse{Store: dto.CustomStoreResponse{ID: *storeID, Name: ""}, Metrics: summary},
        Trend: dto.DashboardTrendResponse{StoreID: *storeID, Range: "12m", Series: trend},
        RecentOrders: dto.DashboardRecentOrdersResponse{StoreID: *storeID, Items: recentOrders},
        LowStock: dto.DashboardLowStockResponse{StoreID: *storeID, Items: lowStock},
        TopProducts: dto.DashboardTopProductsResponse{StoreID: *storeID, Items: topProducts},
        Activity: dto.DashboardActivityResponse{StoreID: *storeID, Items: activity},
    }

    return res, nil
}
