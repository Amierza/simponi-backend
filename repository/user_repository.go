package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	IUserRepository interface {
		// Get
		GetUserByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, bool, error)
		GetLogByStoreID(ctx context.Context, tx *gorm.DB, storeID string) ([]entity.Log, bool, error)
		GetLogByDateRange(ctx context.Context, tx *gorm.DB, startDate, endDate string) ([]entity.Log, bool, error)
		GetAllLogs(ctx context.Context, tx *gorm.DB) ([]entity.Log, bool, error)

		// Create
		CreateLog(ctx context.Context, tx *gorm.DB, log *entity.Log) (*entity.Log, error)
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) GetUserByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	user := new(entity.User)
	err := tx.WithContext(ctx).Where("id = ?", id).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (ur *userRepository) GetLogByStoreID(ctx context.Context, tx *gorm.DB, storeID string) ([]entity.Log, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var logs []entity.Log

	if err := tx.WithContext(ctx).Model(&entity.Log{}).Where("store_id = ?", storeID).Find(&logs).Error; err != nil {
		return nil, false, err
	}

	return logs, true, nil
}

func (ur *userRepository) GetLogByDateRange(ctx context.Context, tx *gorm.DB, startDate, endDate string) ([]entity.Log, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var logs []entity.Log

	if err := tx.WithContext(ctx).Model(&entity.Log{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&logs).Error; err != nil {
		return nil, false, err
	}

	return logs, true, nil
}

func (ur *userRepository) GetAllLogs(ctx context.Context, tx *gorm.DB) ([]entity.Log, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	var logs []entity.Log

	if err := tx.WithContext(ctx).Model(&entity.Log{}).Find(&logs).Error; err != nil {
		return nil, false, err
	}

	return logs, true, nil
}

func (ur *userRepository) CreateLog(ctx context.Context, tx *gorm.DB, log *entity.Log) (*entity.Log, error) {
	if tx == nil {
		tx = ur.db
	}

	if err := tx.WithContext(ctx).Create(log).Error; err != nil {
		return nil, err
	}

	return log, nil
}
