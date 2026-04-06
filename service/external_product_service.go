package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IExternalProductService interface {
		CreateExternalProduct(ctx context.Context, req *dto.CreateExternalProductRequest) (dto.ExternalProductResponse, error)
		GetExternalProducts(ctx context.Context, req *response.PaginationRequest) (dto.ExternalProductPaginationResponse, error)
		GetExternalProductByID(ctx context.Context, externalProductID *uuid.UUID) (dto.ExternalProductResponse, error)
		GetExternalProductsByProductID(ctx context.Context, productID *uuid.UUID) ([]dto.ExternalProductResponse, error)
		GetExternalProductsByStorePlatformID(ctx context.Context, storePlatformID *uuid.UUID) ([]dto.ExternalProductResponse, error)
		UpdateExternalProduct(ctx context.Context, externalProductID *uuid.UUID, req *dto.UpdateExternalProductRequest) (dto.ExternalProductResponse, error)
		DeleteExternalProductByID(ctx context.Context, externalProductID *uuid.UUID) error
	}

	externalProductService struct {
		externalProductRepo repository.IExternalProductRepository
		productRepo         repository.IProductRepository
		logger              *zap.Logger
		jwtService          jwt.IJWT
	}
)

func NewExternalProductService(
	externalProductRepo repository.IExternalProductRepository,
	productRepo repository.IProductRepository,
	logger *zap.Logger,
	jwtService jwt.IJWT,
) *externalProductService {
	return &externalProductService{
		externalProductRepo: externalProductRepo,
		productRepo:         productRepo,
		logger:              logger,
		jwtService:          jwtService,
	}
}

func mapEntityToExternalProductResponse(ep entity.ExternalProduct) dto.ExternalProductResponse {
	return dto.ExternalProductResponse{
		ID:              ep.ID,
		ProductID:       ep.ProductID,
		StorePlatformID: ep.StorePlatformID,
		Price:           ep.Price,
		CreatedAt:       ep.CreatedAt,
		UpdatedAt:       ep.UpdatedAt,
	}
}

func (eps *externalProductService) CreateExternalProduct(ctx context.Context, req *dto.CreateExternalProductRequest) (dto.ExternalProductResponse, error) {
	// validate product exists
	_, found, err := eps.productRepo.GetProductByID(ctx, nil, req.ProductID)
	if err != nil {
		eps.logger.Error("failed to get product by ID", zap.String("productID", req.ProductID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("product not found", zap.String("productID", req.ProductID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	newExternalProduct, err := eps.externalProductRepo.CreateExternalProduct(ctx, nil, &entity.ExternalProduct{
		ID:              uuid.New(),
		ProductID:       req.ProductID,
		StorePlatformID: req.StorePlatformID,
		Price:           req.Price,
	})
	if err != nil {
		eps.logger.Error("failed to create external product", zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to create external product: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to create external product", zap.String("id", newExternalProduct.ID.String()))

	return mapEntityToExternalProductResponse(*newExternalProduct), nil
}

func (eps *externalProductService) GetExternalProducts(ctx context.Context, req *response.PaginationRequest) (dto.ExternalProductPaginationResponse, error) {
	datas, err := eps.externalProductRepo.GetExternalProducts(ctx, nil, req)
	if err != nil {
		eps.logger.Error("failed to get external products", zap.Error(err))
		return dto.ExternalProductPaginationResponse{}, fmt.Errorf("failed to get external products: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to get external products", zap.Int64("count", datas.Count))

	var externalProducts []dto.ExternalProductResponse
	for _, ep := range datas.ExternalProducts {
		externalProducts = append(externalProducts, mapEntityToExternalProductResponse(ep))
	}

	return dto.ExternalProductPaginationResponse{
		Data:               externalProducts,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (eps *externalProductService) GetExternalProductByID(ctx context.Context, externalProductID *uuid.UUID) (dto.ExternalProductResponse, error) {
	externalProduct, found, err := eps.externalProductRepo.GetExternalProductByID(ctx, nil, externalProductID)
	if err != nil {
		eps.logger.Error("failed to get external product by ID", zap.String("externalProductID", externalProductID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get external product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("external product not found", zap.String("externalProductID", externalProductID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("external product not found: %w", dto.ErrNotFound)
	}

	eps.logger.Info("success to get external product by ID", zap.String("id", externalProductID.String()))

	return mapEntityToExternalProductResponse(*externalProduct), nil
}

func (eps *externalProductService) GetExternalProductsByProductID(ctx context.Context, productID *uuid.UUID) ([]dto.ExternalProductResponse, error) {
	externalProducts, err := eps.externalProductRepo.GetExternalProductByProductID(ctx, nil, productID)
	if err != nil {
		eps.logger.Error("failed to get external products by product ID", zap.String("productID", productID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get external products by product ID: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to get external products by product ID", zap.String("productID", productID.String()), zap.Int("count", len(externalProducts)))

	var result []dto.ExternalProductResponse
	for _, ep := range externalProducts {
		result = append(result, mapEntityToExternalProductResponse(ep))
	}

	return result, nil
}

func (eps *externalProductService) GetExternalProductsByStorePlatformID(ctx context.Context, storePlatformID *uuid.UUID) ([]dto.ExternalProductResponse, error) {
	externalProducts, err := eps.externalProductRepo.GetExternalProductByStorePlatformID(ctx, nil, storePlatformID)
	if err != nil {
		eps.logger.Error("failed to get external products by store platform ID", zap.String("storePlatformID", storePlatformID.String()), zap.Error(err))
		return nil, fmt.Errorf("failed to get external products by store platform ID: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to get external products by store platform ID", zap.String("storePlatformID", storePlatformID.String()), zap.Int("count", len(externalProducts)))

	var result []dto.ExternalProductResponse
	for _, ep := range externalProducts {
		result = append(result, mapEntityToExternalProductResponse(ep))
	}

	return result, nil
}

func (eps *externalProductService) UpdateExternalProduct(ctx context.Context, externalProductID *uuid.UUID, req *dto.UpdateExternalProductRequest) (dto.ExternalProductResponse, error) {
	externalProduct, found, err := eps.externalProductRepo.GetExternalProductByID(ctx, nil, externalProductID)
	if err != nil {
		eps.logger.Error("failed to get external product by ID", zap.String("externalProductID", externalProductID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get external product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("external product not found", zap.String("externalProductID", externalProductID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("external product not found: %w", dto.ErrNotFound)
	}

	externalProduct.Price = req.Price

	updatedExternalProduct, err := eps.externalProductRepo.UpdateExternalProduct(ctx, nil, externalProduct)
	if err != nil {
		eps.logger.Error("failed to update external product", zap.String("id", externalProductID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to update external product: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to update external product", zap.String("id", externalProductID.String()))

	return mapEntityToExternalProductResponse(*updatedExternalProduct), nil
}

func (eps *externalProductService) DeleteExternalProductByID(ctx context.Context, externalProductID *uuid.UUID) error {
	_, found, err := eps.externalProductRepo.GetExternalProductByID(ctx, nil, externalProductID)
	if err != nil {
		eps.logger.Error("failed to get external product by ID", zap.String("externalProductID", externalProductID.String()), zap.Error(err))
		return fmt.Errorf("failed to get external product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("external product not found", zap.String("externalProductID", externalProductID.String()))
		return fmt.Errorf("external product not found: %w", dto.ErrNotFound)
	}

	if err := eps.externalProductRepo.DeleteExternalProductByID(ctx, nil, externalProductID); err != nil {
		eps.logger.Error("failed to delete external product by ID", zap.String("externalProductID", externalProductID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete external product by ID: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to delete external product", zap.String("id", externalProductID.String()))

	return nil
}