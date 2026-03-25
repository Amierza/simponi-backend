package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	IUserRepository interface {
		GetUserByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, bool, error)
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

func (ur *userRepository) GetUserByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, bool, error){
	if tx == nil {
		tx = ur.db
	}

	user := new(entity.User)
	err := tx.WithContext(ctx).Where("id = ?", id).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}
