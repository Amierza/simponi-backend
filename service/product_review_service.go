package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type (
	IProductReviewService interface {
		CreateProductReview(ctx context.Context, req *dto.CreateProductReviewRequest) (dto.ProductReviewResponse, error)
		GetProductReviews(ctx context.Context, req *response.PaginationRequest, storeID, productID *uuid.UUID) (dto.ProductReviewPaginationResponse, error)
	}

	productReviewService struct {
		productReviewRepo repository.IProductReviewRepository
		productRepo       repository.IProductRepository
		mlService         IMLService
		logger            *zap.Logger
		jwtService        jwt.IJWT
	}
)

func NewProductReviewService(productReviewRepo repository.IProductReviewRepository, productRepo repository.IProductRepository, mlService IMLService, logger *zap.Logger, jwtService jwt.IJWT) *productReviewService {
	return &productReviewService{
		productReviewRepo: productReviewRepo,
		productRepo:       productRepo,
		mlService:         mlService,
		logger:            logger,
		jwtService:        jwtService,
	}
}

func MapToProductReviewResponse(r entity.ProductReview) dto.ProductReviewResponse {
	tags := []string{}
	if len(r.Tags) > 0 {
		if err := json.Unmarshal(r.Tags, &tags); err != nil {
			tags = []string{}
		}
	}

	productID := uuid.Nil
	if r.ProductID != nil {
		productID = *r.ProductID
	}

	return dto.ProductReviewResponse{
		ID:         r.ID,
		ProductID:  productID,
		ReviewText: r.ReviewText,
		Tags:       tags,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func (prs *productReviewService) CreateProductReview(ctx context.Context, req *dto.CreateProductReviewRequest) (dto.ProductReviewResponse, error) {
	// validate product exists and belongs to the store
	_, found, err := prs.productRepo.GetProductByStoreIDAndProductID(ctx, nil, req.StoreID, req.ProductID)
	if err != nil {
		prs.logger.Error("failed to get product by ID", zap.String("productID", req.ProductID.String()), zap.Error(err))
		return dto.ProductReviewResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		prs.logger.Warn("product not found", zap.String("productID", req.ProductID.String()))
		return dto.ProductReviewResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	// hit ml api to get auto tags
	tags, err := prs.mlService.PredictTags(ctx, req.Text)
	if err != nil {
		prs.logger.Error("failed to predict tags", zap.Error(err))
		return dto.ProductReviewResponse{}, fmt.Errorf("failed to predict tags: %w", err)
	}

	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		prs.logger.Error("failed to marshal tags", zap.Error(err))
		return dto.ProductReviewResponse{}, fmt.Errorf("failed to marshal tags: %w", dto.ErrInternal)
	}

	newReview, err := prs.productReviewRepo.CreateProductReview(ctx, nil, &entity.ProductReview{
		ID:         uuid.New(),
		ProductID:  req.ProductID,
		ReviewText: req.Text,
		Tags:       datatypes.JSON(tagsJSON),
	})
	if err != nil {
		prs.logger.Error("failed to create product review", zap.Error(err))
		return dto.ProductReviewResponse{}, fmt.Errorf("failed to create product review: %w", dto.ErrInternal)
	}

	prs.logger.Info("success to create product review", zap.String("id", newReview.ID.String()))
	return MapToProductReviewResponse(*newReview), nil
}

func (prs *productReviewService) GetProductReviews(ctx context.Context, req *response.PaginationRequest, storeID, productID *uuid.UUID) (dto.ProductReviewPaginationResponse, error) {
	// validate product exists and belongs to the store
	_, found, err := prs.productRepo.GetProductByStoreIDAndProductID(ctx, nil, storeID, productID)
	if err != nil {
		prs.logger.Error("failed to get product by ID", zap.String("productID", productID.String()), zap.Error(err))
		return dto.ProductReviewPaginationResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		prs.logger.Warn("product not found", zap.String("productID", productID.String()))
		return dto.ProductReviewPaginationResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	result, err := prs.productReviewRepo.GetProductReviews(ctx, nil, req, productID)
	if err != nil {
		prs.logger.Error("failed to get product reviews", zap.String("productID", productID.String()), zap.Error(err))
		return dto.ProductReviewPaginationResponse{}, fmt.Errorf("failed to get product reviews: %w", dto.ErrInternal)
	}

	reviews := make([]dto.ProductReviewResponse, 0, len(result.ProductReviews))
	for _, r := range result.ProductReviews {
		reviews = append(reviews, MapToProductReviewResponse(r))
	}

	prs.logger.Info("success to get product reviews", zap.String("productID", productID.String()))
	return dto.ProductReviewPaginationResponse{
		PaginationResponse: result.PaginationResponse,
		Data:               reviews,
	}, nil
}
