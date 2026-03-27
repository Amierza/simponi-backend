package repository

import (
	"context"

	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	ILogRepository interface {
		CreateLog(ctx context.Context, tx *gorm.DB, log *entity.Log) (*entity.Log, error)

		GetLogs(ctx context.Context, tx *gorm.DB) ([]entity.Log, error)
		GetLogByStoreID(ctx context.Context, tx *gorm.DB, storeID string) ([]entity.Log, error)
		GetLogByDateRange(ctx context.Context, tx *gorm.DB, startDate, endDate string) ([]entity.Log, error)
	}

	logRepository struct {
		db *gorm.DB
	}
)

func NewLoggingRepository(db *gorm.DB) *logRepository {
	return &logRepository{
		db: db,
	}
}

func (lr *logRepository) GetLogByStoreID(ctx context.Context, tx *gorm.DB, storeID string) ([]entity.Log, error) {
	if tx == nil {
		tx = lr.db
	}

	var logs []entity.Log

	if err := tx.WithContext(ctx).Model(&entity.Log{}).Where("store_id = ?", storeID).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (lr *logRepository) GetLogByDateRange(ctx context.Context, tx *gorm.DB, startDate, endDate string) ([]entity.Log, error) {
	if tx == nil {
		tx = lr.db
	}

	var logs []entity.Log

	if err := tx.WithContext(ctx).Model(&entity.Log{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (lr *logRepository) GetLogs(ctx context.Context, tx *gorm.DB) ([]entity.Log, error) {
	if tx == nil {
		tx = lr.db
	}

	var logs []entity.Log

	if err := tx.WithContext(ctx).Model(&entity.Log{}).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
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
