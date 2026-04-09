package service

import (
	"context"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"go.uber.org/zap"
)

type (
	IInventoryLogService interface {
		CreateInventoryLog(ctx context.Context, req dto.InventoryLogRequest) (*dto.InventoryLogResponse, error)
		GetInventoryLogs(ctx context.Context, productID string, req response.PaginationRequest) (*dto.InventoryLogPaginationResponse, error)
	}

	inventoryLogService struct {
		inventoryLogRepo repository.IInventoryLoggingRepository
		logger           *zap.Logger
		jwtService       jwt.IJWT
	}
)

func NewInventoryLoggingService(inventoryLogRepo repository.IInventoryLoggingRepository, logger *zap.Logger, jwtService jwt.IJWT) *inventoryLogService {
	return &inventoryLogService{
		inventoryLogRepo: inventoryLogRepo,
		logger:           logger,
		jwtService:       jwtService,
	}
}

func mapInventoryLogsToResponse(inventoryLogs []entity.InventoryLog) []dto.InventoryLogResponse {
	result := make([]dto.InventoryLogResponse, 0, len(inventoryLogs))
	for _, log := range inventoryLogs {
		result = append(result, dto.InventoryLogResponse{
			ID:        log.ID,
			ProductID: log.ProductID,
			Change:    log.Change,
			Source:    log.Source,
			Note:      log.Note,
			CreatedAt: log.CreatedAt,
		})
	}
	return result
}

func (ils *inventoryLogService) CreateInventoryLog(ctx context.Context, req dto.InventoryLogRequest) (*dto.InventoryLogResponse, error) {
	if req.ProductID == nil || req.Change == 0 || req.Source == "" {
		return nil, dto.ErrBadRequest
	}
	log, err := ils.inventoryLogRepo.CreateInventoryLog(ctx, nil, &entity.InventoryLog{
		ProductID: req.ProductID,
		Change:    req.Change,
		Source:    req.Source,
		Note:      req.Note,
	})
	if err != nil {
		ils.logger.Error("Failed to create inventory log", zap.Error(err))
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "foreign key constraint") {
			return nil, dto.ErrBadRequest
		}
		return nil, dto.ErrCreateInventoryLog
	}
	return &dto.InventoryLogResponse{
		ID:        log.ID,
		ProductID: log.ProductID,
		Change:    log.Change,
		Source:    log.Source,
		Note:      log.Note,
		CreatedAt: log.CreatedAt,
	}, nil
}

func (ils *inventoryLogService) GetInventoryLogs(ctx context.Context, productID string, req response.PaginationRequest) (*dto.InventoryLogPaginationResponse, error) {
	datas, err := ils.inventoryLogRepo.GetInventoryLogs(ctx, nil, productID, &req)
	if err != nil {
		ils.logger.Error("Failed to get inventory logs", zap.Error(err))
		return nil, dto.ErrGetInventoryLogs
	}

	return &dto.InventoryLogPaginationResponse{
		PaginationResponse: datas.PaginationResponse,
		Data:               mapInventoryLogsToResponse(datas.InventoryLogs),
	}, nil
}
