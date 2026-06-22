package repository

import (
	"context"
	"math"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IProductReviewRepository interface {
		CreateProductReview(ctx context.Context, tx *gorm.DB, review *entity.ProductReview) (*entity.ProductReview, error)
		GetProductReviews(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, productID *uuid.UUID) (dto.ProductReviewPaginationRepositoryResponse, error)
	}

	productReviewRepository struct {
		db *gorm.DB
	}
)

func NewProductReviewRepository(db *gorm.DB) *productReviewRepository {
	return &productReviewRepository{
		db: db,
	}
}

func (prr *productReviewRepository) CreateProductReview(ctx context.Context, tx *gorm.DB, review *entity.ProductReview) (*entity.ProductReview, error) {
	if tx == nil {
		tx = prr.db
	}

	if err := tx.WithContext(ctx).Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

func (prr *productReviewRepository) GetProductReviews(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, productID *uuid.UUID) (dto.ProductReviewPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = prr.db
	}

	var reviews []entity.ProductReview
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.ProductReview{}).
		Where("product_id = ?", productID)

	if err := query.Count(&count).Error; err != nil {
		return dto.ProductReviewPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&reviews).Error; err != nil {
		return dto.ProductReviewPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.ProductReviewPaginationRepositoryResponse{
		ProductReviews: reviews,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil
}
