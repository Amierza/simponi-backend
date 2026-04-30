package service

import (
	"context"
	"fmt"
	"strings"

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
		GetExternalProducts(ctx context.Context, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ExternalProductPaginationResponse, error)
		GetExternalProductByStoreIDAndExprodID(ctx context.Context, storeID *uuid.UUID, externalProductID *uuid.UUID) (dto.ExternalProductResponse, error)
		GetExternalProductsByStoreIDAndStorePlatformID(ctx context.Context, storeID *uuid.UUID, storePlatformID *uuid.UUID) ([]dto.ExternalProductResponse, error)
		UpdateExternalProductByStoreIDAndExprodID(ctx context.Context, req *dto.UpdateExternalProductRequest) (dto.ExternalProductResponse, error)
		DeleteExternalProductByStoreIDAndExprodID(ctx context.Context, storeID *uuid.UUID, externalProductID *uuid.UUID) error
	}

	externalProductService struct {
		externalProductRepo repository.IExternalProductRepository
		productRepo         repository.IProductRepository
		storeRepo           repository.IStoreRepository
		platformRepo        repository.IPlatformRepository
		storePlatformRepo   repository.IStorePlatformRepository
		logger              *zap.Logger
		jwtService          jwt.IJWT
	}
)

func NewExternalProductService(
	externalProductRepo repository.IExternalProductRepository,
	productRepo repository.IProductRepository,
	storeRepo repository.IStoreRepository,
	platformRepo repository.IPlatformRepository,
	storePlatformRepo repository.IStorePlatformRepository,
	logger *zap.Logger,
	jwtService jwt.IJWT,
) *externalProductService {
	return &externalProductService{
		externalProductRepo: externalProductRepo,
		productRepo:         productRepo,
		storeRepo:           storeRepo,
		platformRepo:        platformRepo,
		storePlatformRepo:   storePlatformRepo,
		logger:              logger,
		jwtService:          jwtService,
	}
}

func mapEntityToExternalProductResponse(ep entity.ExternalProduct) dto.ExternalProductResponse {
	imageURL := ""
	productName := ""
	platformName := ""
	storePlatformName := ""

	if ep.Product != nil {
		productName = ep.Product.Name
		if len(ep.Product.Images) > 0 {
			imageURL = ep.Product.Images[0].ImageURL
		}
	}

	if ep.StorePlatform.Platform != nil {
		platformName = ep.StorePlatform.Platform.Name
	}

	if ep.StorePlatform.Store != nil && ep.StorePlatform.Platform != nil {
		storePlatformName = strings.TrimSpace(ep.StorePlatform.Store.Name + " - " + ep.StorePlatform.Platform.Name)
	}

	return dto.ExternalProductResponse{
		ID:                ep.ID,
		ImageURL:          imageURL,
		ProductName:       productName,
		Platform:          platformName,
		StorePlatformName: storePlatformName,
		Price:             ep.Price,
		CreatedAt:         ep.CreatedAt,
		UpdatedAt:         ep.UpdatedAt,
	}
}

func (eps *externalProductService) CreateExternalProduct(ctx context.Context, req *dto.CreateExternalProductRequest) (dto.ExternalProductResponse, error) {
	_, found, err := eps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, req.StoreID, req.ProductID)
	if err != nil {
		eps.logger.Error("failed to get product by ID", zap.String("productID", req.ProductID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("product not found", zap.String("productID", req.ProductID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	_, found, err = eps.platformRepo.GetPlatformByPlatformID(ctx, nil, req.PlatformID)
	if err != nil {
		eps.logger.Error("failed to get platform by ID", zap.String("platformID", req.PlatformID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get platform ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("platform not found", zap.String("platformID", req.PlatformID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("platform not found: %v", dto.ErrNotFound)
	}

	storePlatform, found, err := eps.storePlatformRepo.GetStorePlatformByStoreIDAndPlatformID(ctx, nil, req.StoreID, req.PlatformID)
	if err != nil {
		eps.logger.Error("failed to get store platform by platform ID", zap.String("platform ID", req.PlatformID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("store platform not found", zap.String("platformID", req.PlatformID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("store platform not found: %w", dto.ErrNotFound)
	}

	newExternalProduct, err := eps.externalProductRepo.CreateExternalProduct(ctx, nil, &entity.ExternalProduct{
		ID:              uuid.New(),
		ProductID:       req.ProductID,
		StorePlatformID: &storePlatform.ID,
		Price:           req.Price,
	})
	if err != nil {
		eps.logger.Error("failed to create external product", zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to create external product: %w", dto.ErrInternal)
	}

	newExternalProduct, _, err = eps.externalProductRepo.GetExternalProductByStoreIDAndExprodID(ctx, nil, req.StoreID, &newExternalProduct.ID)
	if err != nil {
		eps.logger.Error("failed to reload external product after create", zap.String("id", newExternalProduct.ID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to load external product: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to create external product", zap.String("external_product_id", newExternalProduct.ID.String()))

	return mapEntityToExternalProductResponse(*newExternalProduct), nil
}

func (eps *externalProductService) GetExternalProducts(ctx context.Context, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ExternalProductPaginationResponse, error) {
	datas, err := eps.externalProductRepo.GetExternalProducts(ctx, nil, req, storeID)
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

func (eps *externalProductService) GetExternalProductByStoreIDAndExprodID(ctx context.Context, storeID *uuid.UUID, externalProductID *uuid.UUID) (dto.ExternalProductResponse, error) {
	externalProduct, found, err := eps.externalProductRepo.GetExternalProductByStoreIDAndExprodID(ctx, nil, storeID, externalProductID)
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

func (eps *externalProductService) GetExternalProductsByStoreIDAndStorePlatformID(ctx context.Context, storeID *uuid.UUID, storePlatformID *uuid.UUID) ([]dto.ExternalProductResponse, error) {
	externalProducts, err := eps.externalProductRepo.GetExternalProductsByStoreIDAndStorePlatformID(ctx, nil, storeID, storePlatformID)
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

func (eps *externalProductService) UpdateExternalProductByStoreIDAndExprodID(ctx context.Context, req *dto.UpdateExternalProductRequest) (dto.ExternalProductResponse, error) {
	externalProduct, found, err := eps.externalProductRepo.GetExternalProductByStoreIDAndExprodID(ctx, nil, req.StoreID, &req.ID)
	if err != nil {
		eps.logger.Error("failed to get external product by ID", zap.String("externalProductID", req.ID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to get external product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("external product not found", zap.String("externalProductID", req.ID.String()))
		return dto.ExternalProductResponse{}, fmt.Errorf("external product not found: %w", dto.ErrNotFound)
	}

	externalProduct.Price = req.Price

	updatedExternalProduct, err := eps.externalProductRepo.UpdateExternalProductByStoreIDAndExprodID(ctx, nil, externalProduct)
	if err != nil {
		eps.logger.Error("failed to update external product", zap.String("id", req.ID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to update external product: %w", dto.ErrInternal)
	}

	updatedExternalProduct, _, err = eps.externalProductRepo.GetExternalProductByStoreIDAndExprodID(ctx, nil, &updatedExternalProduct.Product.StoreID, &updatedExternalProduct.ID)
	if err != nil {
		eps.logger.Error("failed to reload external product after update", zap.String("id", req.ID.String()), zap.Error(err))
		return dto.ExternalProductResponse{}, fmt.Errorf("failed to load external product: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to update external product", zap.String("id", req.ID.String()))

	return mapEntityToExternalProductResponse(*updatedExternalProduct), nil
}

func (eps *externalProductService) DeleteExternalProductByStoreIDAndExprodID(ctx context.Context, storeID *uuid.UUID, externalProductID *uuid.UUID) error {
	_, found, err := eps.externalProductRepo.GetExternalProductByStoreIDAndExprodID(ctx, nil, storeID, externalProductID)
	if err != nil {
		eps.logger.Error("failed to get external product by ID", zap.String("externalProductID", externalProductID.String()), zap.Error(err))
		return fmt.Errorf("failed to get external product by ID: %w", dto.ErrInternal)
	}
	if !found {
		eps.logger.Warn("external product not found", zap.String("externalProductID", externalProductID.String()))
		return fmt.Errorf("external product not found: %w", dto.ErrNotFound)
	}

	if err := eps.externalProductRepo.DeleteExternalProductByStoreIDAndExprodID(ctx, nil, externalProductID); err != nil {
		eps.logger.Error("failed to delete external product by ID", zap.String("externalProductID", externalProductID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete external product by ID: %w", dto.ErrInternal)
	}

	eps.logger.Info("success to delete external product", zap.String("id", externalProductID.String()))

	return nil
}
