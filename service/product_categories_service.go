package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"go.uber.org/zap"
)

type (
	IProductCategoriesService interface {
		GetProductCategories(ctx context.Context) ([]dto.ProductCategoryResponse, error)
	}

	productCategoriesService struct {
		productCategoriesRepo repository.IProductCategoriesRepository
		logger                *zap.Logger
		jwtService            jwt.IJWT
	}
)

func NewProductCategoriesService(productCategoriesRepo repository.IProductCategoriesRepository, logger *zap.Logger, jwtService jwt.IJWT) *productCategoriesService {
	return &productCategoriesService{
		productCategoriesRepo: productCategoriesRepo,
		logger:                logger,
		jwtService:            jwtService,
	}
}

func (ps *productCategoriesService) GetProductCategories(ctx context.Context) ([]dto.ProductCategoryResponse, error) {
	categories, err := ps.productCategoriesRepo.GetProductCategories(ctx, nil)
	if err != nil {
		ps.logger.Error("failed to get product categories", zap.Error(err))
		return nil, fmt.Errorf("failed to get product categories: %w", dto.ErrInternal)
	}

	var categoryResponses []dto.ProductCategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, dto.ProductCategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	ps.logger.Info("success to get product categories", zap.Int("count", len(categoryResponses)))

	return categoryResponses, nil
}
