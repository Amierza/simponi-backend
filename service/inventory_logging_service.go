package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	IInventoryLogService interface {
		CreateInventoryLog(ctx context.Context, req dto.InventoryLogRequest) (*dto.InventoryLogResponse, error)
		GetInventoryLogs(ctx context.Context, productID string, req response.PaginationRequest) (*dto.InventoryLogPaginationResponse, error)
		GetInventoryLogByID(ctx context.Context, inventoryLogID string) (*dto.InventoryLogResponse, error)
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

func MapToInventoryLogResponse(il entity.InventoryLog) dto.InventoryLogResponse {
	var product *dto.ProductResponse

	if il.Product != nil && il.Product.ID != uuid.Nil {
		product = &dto.ProductResponse{
			ID:   il.Product.ID,
			Name: il.Product.Name,
		}
	}

	return dto.InventoryLogResponse{
		ID:        il.ID,
		Product:   product,
		Change:    il.Change,
		Source:    il.Source,
		Note:      il.Note,
		CreatedAt: il.CreatedAt,
	}
}

func MapInventoryLogsToResponse(ils []entity.InventoryLog) []dto.InventoryLogResponse {
	res := make([]dto.InventoryLogResponse, 0, len(ils))
	for _, il := range ils {
		res = append(res, MapToInventoryLogResponse(il))
	}
	return res
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
		if strings.Contains(errMsg, "foreign key constraint") ||
			strings.Contains(errMsg, "violates foreign key constraint") ||
			strings.Contains(errMsg, "sqlstate 23503") {
			return nil, dto.ErrBadRequest
		}
		return nil, dto.ErrCreateInventoryLog
	}

	logWithProduct, err := ils.inventoryLogRepo.GetInventoryLogByID(ctx, nil, log.ID.String())
	if err != nil {
		ils.logger.Error("Failed to reload inventory log", zap.Error(err))
		return nil, dto.ErrCreateInventoryLog
	}

	res := MapToInventoryLogResponse(*logWithProduct)
	return &res, nil
}

func (ils *inventoryLogService) GetInventoryLogs(ctx context.Context, productID string, req response.PaginationRequest) (*dto.InventoryLogPaginationResponse, error) {
	datas, err := ils.inventoryLogRepo.GetInventoryLogs(ctx, nil, productID, &req)
	if err != nil {
		ils.logger.Error("Failed to get inventory logs", zap.Error(err))
		return nil, dto.ErrGetInventoryLogs
	}

	return &dto.InventoryLogPaginationResponse{
		PaginationResponse: datas.PaginationResponse,
		Data:               MapInventoryLogsToResponse(datas.InventoryLogs),
	}, nil
}

func (ils *inventoryLogService) GetInventoryLogByID(ctx context.Context, inventoryLogID string) (*dto.InventoryLogResponse, error) {
	if strings.TrimSpace(inventoryLogID) == "" {
		return nil, dto.ErrBadRequest
	}

	inventoryLog, err := ils.inventoryLogRepo.GetInventoryLogByID(ctx, nil, inventoryLogID)
	if err != nil {
		ils.logger.Error("Failed to get inventory log by ID", zap.String("inventoryLogID", inventoryLogID), zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}

		return nil, dto.ErrGetInventoryLogs
	}

	res := MapToInventoryLogResponse(*inventoryLog)
	return &res, nil
}
