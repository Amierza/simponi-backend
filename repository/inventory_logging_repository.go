package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"gorm.io/gorm"
)

type (
	IInventoryLoggingRepository interface {
		CreateInventoryLog(ctx context.Context, tx *gorm.DB, log *entity.InventoryLog) (*entity.InventoryLog, error)
		GetInventoryLogs(ctx context.Context, tx *gorm.DB, productID string, req *response.PaginationRequest) (dto.InventoryLogPaginationRepositoryResponse, error)
	}

	inventoryLoggingRepository struct {
		db *gorm.DB
	}
)

func NewInventoryLoggingRepository(db *gorm.DB) *inventoryLoggingRepository {
	return &inventoryLoggingRepository{
		db: db,
	}
}

func (ilr *inventoryLoggingRepository) CreateInventoryLog(ctx context.Context, tx *gorm.DB, log *entity.InventoryLog) (*entity.InventoryLog, error) {
	if tx == nil {
		tx = ilr.db
	}
	if err := tx.WithContext(ctx).Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (ilr *inventoryLoggingRepository) GetInventoryLogs(ctx context.Context, tx *gorm.DB, productID string, req *response.PaginationRequest) (dto.InventoryLogPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = ilr.db
	}

	var inventoryLogs []entity.InventoryLog
	var count int64

	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.InventoryLog{})

	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.InventoryLogPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&inventoryLogs).Error; err != nil {
		return dto.InventoryLogPaginationRepositoryResponse{}, err
	}
	return dto.InventoryLogPaginationRepositoryResponse{
		InventoryLogs: inventoryLogs,
	}, nil
}
