package repository

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IStoreUserRepository interface {
		// CREATE
		CreateStoreUser(ctx context.Context, tx *gorm.DB, storeUser *entity.StoreUser) error

		// READ
		GetStoreUsersByStoreID(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.StoreUserPaginationRepositoryResponse, error)
		GetStoreUserByStoreIDAndUserID(ctx context.Context, tx *gorm.DB, storeID, userID *uuid.UUID) (*entity.StoreUser, bool, error)
		GetStoreUserByUserID(ctx context.Context, tx *gorm.DB, userID *uuid.UUID) (*entity.StoreUser, bool, error)
		// GetStoreUsersByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) ([]*entity.StoreUser, error)

		// UPDATE

		// DELETE
		DeleteStoreUserByStoreIDAndUserID(ctx context.Context, tx *gorm.DB, storeID, userID *uuid.UUID) error
	}

	storeUserRepository struct {
		db *gorm.DB
	}
)

func NewStoreUserRepository(db *gorm.DB) *storeUserRepository {
	return &storeUserRepository{
		db: db,
	}
}

// CREATE
func (sur *storeUserRepository) CreateStoreUser(ctx context.Context, tx *gorm.DB, storeUser *entity.StoreUser) error {
	if tx == nil {
		tx = sur.db
	}

	return tx.WithContext(ctx).Create(storeUser).Error
}

// READ
func (sur *storeUserRepository) GetStoreUsersByStoreID(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.StoreUserPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = sur.db
	}

	var storeUsers []*entity.StoreUser
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.StoreUser{}).
		Preload("User").
		Preload("Store").
		Where("store_id = ?", storeID)

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.
			Joins("JOIN users u ON u.ID = store_users.user_id").
			Where("LOWER(u.name) LIKE ? OR LOWER(u.email) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.StoreUserPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&storeUsers).Error; err != nil {
		return dto.StoreUserPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.StoreUserPaginationRepositoryResponse{
		StoreUsers: storeUsers,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (sur *storeUserRepository) GetStoreUserByStoreIDAndUserID(ctx context.Context, tx *gorm.DB, storeID, userID *uuid.UUID) (*entity.StoreUser, bool, error) {
	if tx == nil {
		tx = sur.db
	}

	var storeUser *entity.StoreUser
	err := tx.WithContext(ctx).
		Model(&entity.StoreUser{}).
		Preload("Store").
		Preload("User").
		Preload("User.Role").
		Where("store_id = ?", storeID).
		Where("user_id = ?", userID).
		First(&storeUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return storeUser, true, nil
}
func (sur *storeUserRepository) GetStoreUserByUserID(ctx context.Context, tx *gorm.DB, userID *uuid.UUID) (*entity.StoreUser, bool, error) {
	if tx == nil {
		tx = sur.db
	}

	var storeUser *entity.StoreUser
	err := tx.WithContext(ctx).
		Model(&entity.StoreUser{}).
		Preload("Store").
		Preload("User").
		Where("user_id = ?", userID).
		First(&storeUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return storeUser, true, nil
}

// func (sur *storeUserRepository) GetStoreUsersByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) ([]*entity.StoreUser, error) {
// 	if tx == nil {
// 		tx = sur.db
// 	}

// 	var storeUsers []*entity.StoreUser

// 	if err := tx.WithContext(ctx).
// 		Model(&entity.StoreUser{}).
// 		Preload("User").
// 		Preload("Store").
// 		Where("store_id = ?", storeID).
// 		Find(&storeUsers).Error; err != nil {
// 		return nil, err
// 	}

// 	return storeUsers, nil
// }

// UPDATE

// DELETE
func (sur *storeUserRepository) DeleteStoreUserByStoreIDAndUserID(ctx context.Context, tx *gorm.DB, storeID, userID *uuid.UUID) error {
	if tx == nil {
		tx = sur.db
	}

	return tx.WithContext(ctx).Where("store_id = ? AND user_id = ?", &storeID, &userID).Delete(&entity.StoreUser{}).Error
}
