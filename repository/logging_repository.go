package repository

import (
	"context"
	"math"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"gorm.io/gorm"
)

type (
	ILogRepository interface {
		CreateLog(ctx context.Context, tx *gorm.DB, log *entity.Log) (*entity.Log, error)
		GetLogs(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.LogPaginationRepositoryResponse, error)
		GetLogByStoreID(ctx context.Context, tx *gorm.DB, storeID string, req *response.PaginationRequest) (dto.LogPaginationRepositoryResponse, error)
		GetLogByDateRange(ctx context.Context, tx *gorm.DB, startDate, endDate string, req *response.PaginationRequest) (dto.LogPaginationRepositoryResponse, error)
	}

	logRepository struct {
		db *gorm.DB
	}
)

func NewLogRepository(db *gorm.DB) *logRepository {
	return &logRepository{
		db: db,
	}
}

func (lr *logRepository) GetLogs(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.LogPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = lr.db
	}

	var logs []entity.Log
	var count int64

	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Log{})

	if err := query.Count(&count).Error; err != nil {
		return dto.LogPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&logs).Error; err != nil {
		return dto.LogPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.LogPaginationRepositoryResponse{
		Logs: logs,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (lr *logRepository) GetLogByStoreID(ctx context.Context, tx *gorm.DB, storeID string, req *response.PaginationRequest) (dto.LogPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = lr.db
	}

	var logs []entity.Log
	var count int64

	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Log{}).Where("store_id = ?", storeID)

	if err := query.Count(&count).Error; err != nil {
		return dto.LogPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&logs).Error; err != nil {
		return dto.LogPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.LogPaginationRepositoryResponse{
		Logs: logs,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (lr *logRepository) GetLogByDateRange(ctx context.Context, tx *gorm.DB, startDate, endDate string, req *response.PaginationRequest) (dto.LogPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = lr.db
	}

	var logs []entity.Log
	var count int64

	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Log{}).Where("created_at BETWEEN ? AND ?", startDate, endDate)

	if err := query.Count(&count).Error; err != nil {
		return dto.LogPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&logs).Error; err != nil {
		return dto.LogPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.LogPaginationRepositoryResponse{
		Logs: logs,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}

func (lr *logRepository) CreateLog(ctx context.Context, tx *gorm.DB, log *entity.Log) (*entity.Log, error) {
	if tx == nil {
		tx = lr.db
	}

	if err := tx.WithContext(ctx).Create(log).Error; err != nil {
		return nil, err
	}

	return log, nil
}
