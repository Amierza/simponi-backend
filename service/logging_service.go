package service

import (
	"context"
	"fmt"
	"strings"
	"time"

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
		GetLogs(ctx context.Context, storeID, startDate, endDate string, req response.PaginationRequest) (*dto.LogPaginationResponse, error)
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
	if req.StoreID == nil || strings.TrimSpace(req.Action) == "" || strings.TrimSpace(req.Message) == "" {
		return nil, dto.ErrBadRequest
	}

	log, err := ls.logRepo.CreateLog(ctx, nil, &entity.Log{
		StoreID: req.StoreID,
		Action:  req.Action,
		Message: req.Message,
	})
	if err != nil {
		ls.logger.Error("failed to create log", zap.Error(err))
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "violates foreign key constraint") || strings.Contains(errMsg, "sqlstate 23503") {
			return nil, dto.ErrBadRequest
		}
		return nil, dto.ErrCreateLog
	}

	return &dto.LogResponse{
		ID:        log.ID,
		StoreID:   log.StoreID,
		Action:    log.Action,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (ls *logService) GetLogs(ctx context.Context, storeID, startDate, endDate string, req response.PaginationRequest) (*dto.LogPaginationResponse, error) {
	if strings.TrimSpace(startDate) != "" {
		if _, err := time.Parse(time.DateOnly, startDate); err != nil {
			return nil, fmt.Errorf("%w: start_date must be YYYY-MM-DD", dto.ErrBadRequest)
		}
	}

	if strings.TrimSpace(endDate) != "" {
		if _, err := time.Parse(time.DateOnly, endDate); err != nil {
			return nil, fmt.Errorf("%w: end_date must be YYYY-MM-DD", dto.ErrBadRequest)
		}
	}

	if strings.TrimSpace(startDate) != "" && strings.TrimSpace(endDate) != "" {
		start, _ := time.Parse(time.DateOnly, startDate)
		end, _ := time.Parse(time.DateOnly, endDate)
		if start.After(end) {
			return nil, fmt.Errorf("%w: start_date must be less than or equal to end_date", dto.ErrBadRequest)
		}
	}

	datas, err := ls.logRepo.GetLogs(ctx, nil, strings.TrimSpace(storeID), strings.TrimSpace(startDate), strings.TrimSpace(endDate), &req)
	if err != nil {
		ls.logger.Error("failed to get logs", zap.Error(err))
		return nil, dto.ErrGetLogs
	}

	return &dto.LogPaginationResponse{
		PaginationResponse: datas.PaginationResponse,
		Data:               mapLogsToResponse(datas.Logs),
	}, nil
}
