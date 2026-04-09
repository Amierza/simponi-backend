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
		CreateLog(ctx context.Context, req dto.LogRequest) (*dto.LogResponse, error)
		GetLogs(ctx context.Context, req response.PaginationRequest) (*dto.LogPaginationResponse, error)
		GetLogsByStoreID(ctx context.Context, storeID string, req response.PaginationRequest) (*dto.LogPaginationResponse, error)
		GetLogsByDateRange(ctx context.Context, startDate, endDate string, req response.PaginationRequest) (*dto.LogPaginationResponse, error)
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

func mapLogsToResponse(logs []entity.Log) []dto.LogResponse {
	result := make([]dto.LogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, dto.LogResponse{
			ID:        log.ID,
			StoreID:   log.StoreID,
			Action:    log.Action,
			Message:   log.Message,
			CreatedAt: log.CreatedAt,
		})
	}
	return result
}

func (ls *logService) CreateLog(ctx context.Context, req dto.LogRequest) (*dto.LogResponse, error) {
	log, err := ls.logRepo.CreateLog(ctx, nil, &entity.Log{
		StoreID: req.StoreID,
		Action:  req.Action,
		Message: req.Message,
	})
	if err != nil {
		ls.logger.Error("failed to create log", zap.Error(err))
		return nil, err
	}

	return &dto.LogResponse{
		ID:        log.ID,
		StoreID:   log.StoreID,
		Action:    log.Action,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (ls *logService) GetLogs(ctx context.Context, req response.PaginationRequest) (*dto.LogPaginationResponse, error) {
	datas, err := ls.logRepo.GetLogs(ctx, nil, &req)
	if err != nil {
		ls.logger.Error("failed to get all logs", zap.Error(err))
		return nil, dto.ErrGetLogs
	}

	return &dto.LogPaginationResponse{
		PaginationResponse: datas.PaginationResponse,
		Data:               mapLogsToResponse(datas.Logs),
	}, nil
}

func (ls *logService) GetLogsByStoreID(ctx context.Context, storeID string, req response.PaginationRequest) (*dto.LogPaginationResponse, error) {
	datas, err := ls.logRepo.GetLogByStoreID(ctx, nil, storeID, &req)
	if err != nil {
		ls.logger.Error("failed to get logs by store ID", zap.Error(err))
		return nil, dto.ErrGetLogsByStoreID
	}

	return &dto.LogPaginationResponse{
		PaginationResponse: datas.PaginationResponse,
		Data:               mapLogsToResponse(datas.Logs),
	}, nil
}

func (ls *logService) GetLogsByDateRange(ctx context.Context, startDate, endDate string, req response.PaginationRequest) (*dto.LogPaginationResponse, error) {
	datas, err := ls.logRepo.GetLogByDateRange(ctx, nil, startDate, endDate, &req)
	if err != nil {
		ls.logger.Error("failed to get logs by date range", zap.Error(err))
		return nil, dto.ErrGetLogsByDateRange
	}

	return &dto.LogPaginationResponse{
		PaginationResponse: datas.PaginationResponse,
		Data:               mapLogsToResponse(datas.Logs),
	}, nil
}
