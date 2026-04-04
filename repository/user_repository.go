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
	IUserRepository interface {
		// CREATE
		CreateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error

		// GET
		GetUsers(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.UserPaginationRepositoryResponse, error)
		GetUserByID(ctx context.Context, tx *gorm.DB, userID *uuid.UUID) (*entity.User, bool, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email *string) (*entity.User, bool, error)

		// UPDATE
		UpdateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error

		// DELETE
		DeleteUserByID(ctx context.Context, tx *gorm.DB, userID *uuid.UUID) error
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

// CREATE
func (ur *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Create(&user).Error
}

// GET
func (ur *userRepository) GetUsers(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.UserPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = ur.db
	}

	var users []*entity.User
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.User{}).
		Preload("Role").
		Preload("Stores")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.UserPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&users).Error; err != nil {
		return dto.UserPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.UserPaginationRepositoryResponse{
		Users: users,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (ur *userRepository) GetUserByID(ctx context.Context, tx *gorm.DB, id *uuid.UUID) (*entity.User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	user := new(entity.User)
	err := tx.WithContext(ctx).
		Preload("Role").
		Preload("Stores").
		Where("id = ?", id).
		Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}
func (ur *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email *string) (*entity.User, bool, error) {
	if tx == nil {
		tx = ur.db
	}

	user := new(entity.User)
	err := tx.WithContext(ctx).
		Preload("Role").
		Preload("Stores").Where("email = ?", email).
		Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}

// UPDATE
func (ur *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// DELETE
func (ur *userRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, userID *uuid.UUID) error {
	if tx == nil {
		tx = ur.db
	}

	return tx.WithContext(ctx).Where("id = ?", &userID).Delete(&entity.User{}).Error
}
