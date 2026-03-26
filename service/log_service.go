package service

import (
	"context"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"go.uber.org/zap"
)

type (
	ILogService interface {
		CreateLog(ctx context.Context, req dto.LoggingRequest) (dto.LoggingResponse, error)
		GetLogs(ctx context.Context, req dto.PaginationRequest) (dto.LoggingPaginationResponse, error)
		GetLogsByStoreID(ctx context.Context, storeID string, req dto.PaginationRequest) (dto.LoggingPaginationResponse, error)
		GetLogsByDateRange(ctx context.Context, startDate, endDate string, req dto.PaginationRequest) (dto.LoggingPaginationResponse, error)
	}

	logService struct {
		userRepo repository.IUserRepository
		logger   *zap.Logger
		jwt      jwt.IJWT
	}
)

func NewLogService(userRepo repository.IUserRepository, logger *zap.Logger, jwt jwt.IJWT) *logService {
	return &logService{
		userRepo: userRepo,
		logger:   logger,
		jwt:      jwt,
	}
}

func (ls *logService) CreateLog(ctx context.Context, req dto.LoggingRequest) (dto.LoggingResponse, error) {
	log, err := ls.userRepo.CreateLog(ctx, nil, &entity.Log{
		StoreID: req.StoreID,
		Action:  req.Action,
		Message: req.Message,
	})

	if err != nil {
		ls.logger.Error("failed to create log", zap.Error(err))
		return dto.LoggingResponse{}, err
	}

	return dto.LoggingResponse{
		ID:        log.ID,
		StoreID:   log.StoreID,
		Action:    log.Action,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (ls *logService) GetLogsByStoreID(ctx context.Context, storeID string, req dto.PaginationRequest) ([]dto.LoggingResponse, error) {
	logs, found, err := ls.userRepo.GetLogByStoreID(ctx, nil, storeID)
	if err != nil {
		ls.logger.Error("failed to get logs by store ID", zap.Error(err))
		return nil, err
	}

	if !found {
		ls.logger.Warn("logs not found for store ID", zap.String("storeID", storeID))
		return nil, dto.ErrNotFound
	}

	var paginatedLogs []dto.LoggingResponse
	for _, log := range logs {
		paginatedLogs = append(paginatedLogs, dto.LoggingResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	return paginatedLogs, nil
}

func (ls *logService) GetLogsByDateRange(ctx context.Context, startDate, endDate string, req dto.PaginationRequest) ([]dto.LoggingResponse, error) {
	logs, found, err := ls.userRepo.GetLogByDateRange(ctx, nil, startDate, endDate)
	if err != nil {
		ls.logger.Error("failed to get logs by date range", zap.Error(err))
		return nil, err
	}

	if !found {
		ls.logger.Warn("logs not found for date range", zap.String("startDate", startDate), zap.String("endDate", endDate))
		return nil, dto.ErrNotFound
	}

	var paginatedLogs []dto.LoggingResponse
	for _, log := range logs {
		paginatedLogs = append(paginatedLogs, dto.LoggingResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	return paginatedLogs, nil
}

func (ls *logService) GetAllLogs(ctx context.Context, req dto.PaginationRequest) ([]dto.LoggingResponse, error) {
	logs, found, err := ls.userRepo.GetAllLogs(ctx, nil)

	if err != nil {
		ls.logger.Error("failed to get all logs", zap.Error(err))
		return nil, err
	}

	if !found {
		ls.logger.Warn("logs not found")
		return nil, dto.ErrNotFound
	}

	var paginatedLogs []dto.LoggingResponse
	for _, log := range logs {
		paginatedLogs = append(paginatedLogs, dto.LoggingResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	return paginatedLogs, nil
}
