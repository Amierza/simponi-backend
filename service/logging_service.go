package service

import (
	"context"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"go.uber.org/zap"
)

type (
	ILogService interface {
		CreateLog(ctx context.Context, req dto.LogRequest) (dto.LogResponse, error)
		GetLogs(ctx context.Context, req response.PaginationRequest) (dto.LogPaginationResponse, error)
		GetLogsByStoreID(ctx context.Context, storeID string, req response.PaginationRequest) (dto.LogPaginationResponse, error)
		GetLogsByDateRange(ctx context.Context, startDate, endDate string, req response.PaginationRequest) (dto.LogPaginationResponse, error)
	}

	logService struct {
		logRepo    repository.ILogRepository
		logger     *zap.Logger
		jwtService jwt.IJWT
	}
)

func NewLogService(logRepo repository.ILogRepository, logger *zap.Logger, jwtService jwt.IJWT) *logService {
	return &logService{
		logRepo:    logRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}

func (ls *logService) CreateLog(ctx context.Context, req dto.LogRequest) (dto.LogResponse, error) {
	log, err := ls.logRepo.CreateLog(ctx, nil, &entity.Log{
		StoreID: req.StoreID,
		Action:  req.Action,
		Message: req.Message,
	})

	if err != nil {
		ls.logger.Error("failed to create log", zap.Error(err))
		return dto.LogResponse{}, err
	}

	return dto.LogResponse{
		ID:        log.ID,
		StoreID:   log.StoreID,
		Action:    log.Action,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (ls *logService) GetLogsByStoreID(ctx context.Context, storeID string, req response.PaginationRequest) (dto.LogPaginationResponse, error) {
	logs, err := ls.logRepo.GetLogByStoreID(ctx, nil, storeID)
	if err != nil {
		ls.logger.Error("failed to get logs by store ID", zap.Error(err))
		return dto.LogPaginationResponse{}, dto.ErrGetLogsByStoreID
	}

	if len(logs) == 0 {
		ls.logger.Warn("logs not found for store ID", zap.String("storeID", storeID))
		return dto.LogPaginationResponse{}, dto.ErrNotFound
	}

	var paginatedLogs []dto.LogResponse
	for _, log := range logs {
		paginatedLogs = append(paginatedLogs, dto.LogResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	count := int64(len(paginatedLogs))
	perPage := int64(req.PerPage)
	maxPage := (count + perPage - 1) / perPage

	return dto.LogPaginationResponse{
		Data: paginatedLogs,
		Pagination: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (ls *logService) GetLogsByDateRange(ctx context.Context, startDate, endDate string, req response.PaginationRequest) (dto.LogPaginationResponse, error) {
	logs, err := ls.logRepo.GetLogByDateRange(ctx, nil, startDate, endDate)
	if err != nil {
		ls.logger.Error("failed to get logs by date range", zap.Error(err))
		return dto.LogPaginationResponse{}, dto.ErrGetLogsByDateRange
	}

	if len(logs) == 0 {
		ls.logger.Warn("logs not found for date range", zap.String("startDate", startDate), zap.String("endDate", endDate))
		return dto.LogPaginationResponse{}, dto.ErrNotFound
	}

	var paginatedLogs []dto.LogResponse
	for _, log := range logs {
		paginatedLogs = append(paginatedLogs, dto.LogResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	count := int64(len(paginatedLogs))
	perPage := int64(req.PerPage)
	maxPage := (count + perPage - 1) / perPage

	return dto.LogPaginationResponse{
		Data: paginatedLogs,
		Pagination: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (ls *logService) GetLogs(ctx context.Context, req response.PaginationRequest) (dto.LogPaginationResponse, error) {
	logs, err := ls.logRepo.GetLogs(ctx, nil)

	if err != nil {
		ls.logger.Error("failed to get all logs", zap.Error(err))
		return dto.LogPaginationResponse{}, dto.ErrGetLogs
	}

	if len(logs) == 0 {
		ls.logger.Warn("logs not found")
		return dto.LogPaginationResponse{}, dto.ErrNotFound
	}

	var paginatedLogs []dto.LogResponse
	for _, log := range logs {
		paginatedLogs = append(paginatedLogs, dto.LogResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}

	count := int64(len(paginatedLogs))
	perPage := int64(req.PerPage)
	maxPage := (count + perPage - 1) / perPage

	return dto.LogPaginationResponse{
		Data: paginatedLogs,
		Pagination: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}
