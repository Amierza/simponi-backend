package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"

	"go.uber.org/zap"
)

type (
	IDashboardService interface {
		GetDashboardData(ctx context.Context) (*dto.DashboardSuperadminResponse, error)
	}

	dashboardService struct {
		dashboardRepo repository.IDashboardRepository
		logger        *zap.Logger
		jwtService    jwt.IJWT
	}
)

func NewDashboardService(dashboardRepo repository.IDashboardRepository, logger *zap.Logger, jwtService jwt.IJWT) *dashboardService {
	return &dashboardService{
		dashboardRepo: dashboardRepo,
		logger:        logger,
		jwtService:    jwtService,
	}
}

func (ds *dashboardService) GetDashboardData(ctx context.Context) (*dto.DashboardSuperadminResponse, error) {
	response, err := ds.dashboardRepo.GetDashboardData(ctx, nil)
	if err != nil {
		ds.logger.Error("failed to get dashboard data", zap.Error(err))
		return nil, fmt.Errorf("failed to get dashboard data: %w", dto.ErrInternal)
	}

	return &response, nil
}
