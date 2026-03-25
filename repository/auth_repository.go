package repository

import (
	"context"
	"errors"

	"github.com/Amierza/simponi-backend/entity"
	"gorm.io/gorm"
)

type (
	IAuthRepository interface {
		GetUserByEmail(ctx context.Context, tx *gorm.DB, identifier *string) (*entity.User, bool, error)
	}

	authRepository struct {
		db *gorm.DB
	}
)

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, identifier *string) (*entity.User, bool, error){
	if tx == nil{
		tx = ar.db
	}

	if identifier == nil {
		return nil, false, nil
	}

	user := new(entity.User)
	err := tx.WithContext(ctx).Where("email = ?", *identifier).Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, false, nil
	}
	if err != nil{
		return nil, false, err
	}

	return user, true, nil
}

