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
	IProductService interface {
		CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (dto.ProductResponse, error)
		GetProducts(ctx context.Context, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ProductPaginationResponse, error)
		GetProductStats(ctx context.Context, storeID *uuid.UUID) (dto.ProductStatsResponse, error)
		GetProductByStoreIDAndProductID(ctx context.Context, storeID *uuid.UUID, productID *uuid.UUID) (dto.ProductResponse, error)
		UpdateProductByStoreIDAndProductID(ctx context.Context, req *dto.UpdateProductRequest) (dto.ProductResponse, error)
		UpdateStockByStoreIDAndProductID(ctx context.Context, req *dto.UpdateStockRequest) error
		DeleteProductByStoreIDAndProductID(ctx context.Context, storeID *uuid.UUID, productID *uuid.UUID) error
	}

	productService struct {
		productRepo        repository.IProductRepository
		storeRepo          repository.IStoreRepository
		inventoryLogService IInventoryLogService
		logger             *zap.Logger
		jwtService         jwt.IJWT
	}
)

func NewProductService(productRepo repository.IProductRepository, storeRepo repository.IStoreRepository, inventoryLogService IInventoryLogService, logger *zap.Logger, jwtService jwt.IJWT) *productService {
	return &productService{
		productRepo:         productRepo,
		storeRepo:           storeRepo,
		inventoryLogService: inventoryLogService,
		logger:              logger,
		jwtService:          jwtService,
	}
}

func (ps *productService) createInventoryLog(ctx context.Context, productID *uuid.UUID, change int, source, note string) {
	if ps.inventoryLogService == nil || productID == nil || change == 0 {
		return
	}

	if _, err := ps.inventoryLogService.CreateInventoryLog(ctx, dto.InventoryLogRequest{
		ProductID: productID,
		Change:    change,
		Source:    source,
		Note:      note,
	}); err != nil {
		ps.logger.Warn("failed to create inventory log", zap.String("productID", productID.String()), zap.Error(err))
	}
}

func getProductStatus(p entity.Product) string {
	if p.Stock == 0 {
		return "Out of stock"
	}
	if p.Stock <= 10 {
		return "Low Stock"
	}
	if len(p.ExternalProducts) == 0 {
		return "Unmapped"
	}
	return "Mapped"
}

func MapToProductStoreResponse(p entity.Product) *dto.ProductStoreResponse {
	return &dto.ProductStoreResponse{
		ID:   p.Store.ID,
		Name: p.Store.Name,
	}
}

func MapToProductCategoryResponse(p entity.Product) *dto.ProductCategoryResponse {
	if p.CategoryID == nil {
		return nil
	}
	return &dto.ProductCategoryResponse{
		ID:   p.Category.ID,
		Name: p.Category.Name,
	}
}

func MapToProductImageResponse(p entity.Product) []dto.ProductImageResponse {
	var images []dto.ProductImageResponse

	for _, img := range p.Images {
		images = append(images, dto.ProductImageResponse{
			ID:       img.ID,
			ImageURL: img.ImageURL,
		})
	}
	return images
}

func MapToExternalProductResponse(p entity.Product) []dto.ExternalProductResponse {
	var externalProducts []dto.ExternalProductResponse

	primaryImageURL := ""
	if len(p.Images) > 0 {
		primaryImageURL = p.Images[0].ImageURL
	}

	for _, ep := range p.ExternalProducts {
		platformName := ""
		storePlatformName := ""
		if ep.StorePlatform.Platform != nil {
			platformName = ep.StorePlatform.Platform.Name
		}

		if ep.StorePlatform.Store != nil && ep.StorePlatform.Platform != nil {
			storePlatformName = strings.TrimSpace(ep.StorePlatform.Store.Name + " - " + ep.StorePlatform.Platform.Name)
		}

		externalProducts = append(externalProducts, dto.ExternalProductResponse{
			ID:                ep.ID,
			ImageURL:          primaryImageURL,
			ProductName:       p.Name,
			Platform:          platformName,
			StorePlatformName: storePlatformName,
			Price:             ep.Price,
			CreatedAt:         ep.CreatedAt,
			UpdatedAt:         ep.UpdatedAt,
		})
	}
	return externalProducts
}

func MapToProductListResponse(p entity.Product) dto.ProductListResponse {
	return dto.ProductListResponse{
		ID:               p.ID,
		Name:             p.Name,
		SKU:              p.SKU,
		Stock:            p.Stock,
		Store:            MapToProductStoreResponse(p),
		Category:         MapToProductCategoryResponse(p),
		Images:           MapToProductImageResponse(p),
		ExternalProducts: MapToExternalProductResponse(p),
		Status:           getProductStatus(p),
		CreatedAt:        p.CreatedAt,
	}
}

func MapToProductResponse(p entity.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		SKU:              p.SKU,
		Stock:            p.Stock,
		Store:            MapToProductStoreResponse(p),
		Category:         MapToProductCategoryResponse(p),
		Images:           MapToProductImageResponse(p),
		ExternalProducts: MapToExternalProductResponse(p),
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func (ps *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (dto.ProductResponse, error) {
	_, found, err := ps.productRepo.GetProductBySKUAndStoreID(ctx, nil, req.SKU, req.StoreID)
	if err != nil {
		ps.logger.Error("failed to get product by sku", zap.String("sku", req.SKU), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to get product by SKU: %w", dto.ErrInternal)
	}
	if found {
		ps.logger.Warn("product SKU already exists", zap.String("sku", req.SKU))
		return dto.ProductResponse{}, fmt.Errorf("product already exists: %w", dto.ErrAlreadyExists)
	}

	newProduct, err := ps.productRepo.CreateProduct(ctx, nil, &entity.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		StoreID:     *req.StoreID,
	})
	if err != nil {
		ps.logger.Error("failed to create product", zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to create product: %w", dto.ErrInternal)
	}

	if len(req.Images) > 0 {
		for _, imageURL := range req.Images {
			trimmedURL := strings.TrimSpace(imageURL)
			if trimmedURL == "" {
				continue
			}

			_, err := ps.productRepo.CreateProductImage(ctx, nil, &entity.ProductImage{
				ID:        uuid.New(),
				ImageURL:  trimmedURL,
				ProductID: &newProduct.ID,
			})
			if err != nil {
				ps.logger.Error("failed to create product image", zap.String("productID", newProduct.ID.String()), zap.String("imageURL", trimmedURL), zap.Error(err))
				return dto.ProductResponse{}, fmt.Errorf("failed to create product image: %w", dto.ErrInternal)
			}
		}
	}

	newProduct, _, err = ps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, &newProduct.StoreID, &newProduct.ID)
	if err != nil {
		ps.logger.Error("failed to reload product after image attach", zap.String("productID", newProduct.ID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to load product: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to create product", zap.String("id", newProduct.ID.String()))

	ps.createInventoryLog(ctx, &newProduct.ID, newProduct.Stock, "product", "create product")

	return MapToProductResponse(*newProduct), nil
}

func (ps *productService) GetProducts(ctx context.Context, req *response.PaginationRequest, storeID *uuid.UUID) (dto.ProductPaginationResponse, error) {
	datas, err := ps.productRepo.GetProducts(ctx, nil, req, storeID)
	if err != nil {
		ps.logger.Error("failed to get products", zap.Error(err))
		return dto.ProductPaginationResponse{}, fmt.Errorf("failed to get products: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to get products", zap.Int64("count", datas.Count))

	var productList []dto.ProductListResponse
	for _, p := range datas.Products {
		productList = append(productList, MapToProductListResponse(p))
	}

	return dto.ProductPaginationResponse{
		Data:               productList,
		PaginationResponse: datas.PaginationResponse,
	}, nil
}

func (ps *productService) GetProductStats(ctx context.Context, storeID *uuid.UUID) (dto.ProductStatsResponse, error) {
	stats, err := ps.productRepo.GetProductStats(ctx, nil, storeID)
	if err != nil {
		ps.logger.Error("failed to get product stats", zap.Error(err))
		return dto.ProductStatsResponse{}, fmt.Errorf("failed to get product stats: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to get product stats")

	return stats, nil
}

func (ps *productService) GetProductByStoreIDAndProductID(ctx context.Context, storeID *uuid.UUID, productID *uuid.UUID) (dto.ProductResponse, error) {
	product, found, err := ps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, storeID, productID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", productID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		ps.logger.Warn("product not found", zap.String("productID", productID.String()))
		return dto.ProductResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	ps.logger.Info("success to get product by ID", zap.String("id", productID.String()))

	return MapToProductResponse(*product), nil
}

func (ps *productService) UpdateProductByStoreIDAndProductID(ctx context.Context, req *dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, found, err := ps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, req.StoreID, &req.ID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", req.ID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}

	if !found {
		ps.logger.Warn("product not found", zap.String("productID", req.ID.String()))
		return dto.ProductResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	if req.SKU != "" && req.SKU != product.SKU {
		_, found, err := ps.productRepo.GetProductBySKUAndStoreID(ctx, nil, req.SKU, req.StoreID)
		if err != nil {
			ps.logger.Error("failed to get product by SKU", zap.String("sku", req.SKU), zap.Error(err))
			return dto.ProductResponse{}, fmt.Errorf("failed to get product by SKU: %w", dto.ErrInternal)
		}

		if found {
			ps.logger.Warn("product SKU already exists", zap.String("sku", req.SKU))
			return dto.ProductResponse{}, fmt.Errorf("product SKU already exists: %w", dto.ErrAlreadyExists)
		}
		product.SKU = req.SKU
	}

	product.Name = req.Name
	if req.Description != nil {
		product.Description = *req.Description
	}

	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}

	oldStock := product.Stock
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	updatedProduct, err := ps.productRepo.UpdateProductByStoreIDAndProductID(ctx, nil, product)

	if err != nil {
		ps.logger.Error("failed to update product", zap.String("id", req.ID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to update product: %w", dto.ErrInternal)
	}

	if len(req.Images) > 0 {
		for _, imageURL := range req.Images {
			trimmedURL := strings.TrimSpace(imageURL)
			if trimmedURL == "" {
				continue
			}

			_, err := ps.productRepo.CreateProductImage(ctx, nil, &entity.ProductImage{
				ID:        uuid.New(),
				ImageURL:  trimmedURL,
				ProductID: &updatedProduct.ID,
			})
			if err != nil {
				ps.logger.Error("failed to create product image", zap.String("productID", updatedProduct.ID.String()), zap.String("imageURL", trimmedURL), zap.Error(err))
				return dto.ProductResponse{}, fmt.Errorf("failed to create product image: %w", dto.ErrInternal)
			}
		}
	}

	updatedProduct, _, err = ps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, &updatedProduct.StoreID, &updatedProduct.ID)
	if err != nil {
		ps.logger.Error("failed to reload product after image attach", zap.String("productID", updatedProduct.ID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to load product: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to create product", zap.String("id", updatedProduct.ID.String()))

	ps.createInventoryLog(ctx, &updatedProduct.ID, updatedProduct.Stock-oldStock, "product", "update product")

	return MapToProductResponse(*updatedProduct), nil
}

func (ps *productService) UpdateStockByStoreIDAndProductID(ctx context.Context, req *dto.UpdateStockRequest) error {
	_, found, err := ps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, req.StoreID, &req.ID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", req.ID.String()), zap.Error(err))
		return fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}

	if !found {
		ps.logger.Warn("product not found", zap.String("productID", req.ID.String()))
		return fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	if err := ps.productRepo.UpdateStockByStoreIDAndProductID(ctx, nil, req.StoreID, &req.ID, req.Change); err != nil {
		ps.logger.Error("failed to update stock", zap.String("productID", req.ID.String()), zap.Error(err))
		return dto.ErrInternal
	}

	ps.logger.Info("success to update stock", zap.String("categoryID", req.ID.String()), zap.Int("change", req.Change))

	return nil
}

func (ps *productService) DeleteProductByStoreIDAndProductID(ctx context.Context, storeID *uuid.UUID, productID *uuid.UUID) error {
	product, found, err := ps.productRepo.GetProductByStoreIDAndProductID(ctx, nil, storeID, productID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", productID.String()), zap.Error(err))
		return fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}
	if !found {
		ps.logger.Warn("product not found", zap.String("productID", productID.String()))
		return fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	ps.createInventoryLog(ctx, &product.ID, -product.Stock, "product", "delete product")

	if err := ps.productRepo.DeleteProductByStoreIDAndProductID(ctx, nil, storeID, productID); err != nil {
		ps.logger.Error("failed to delete product", zap.String("productID", productID.String()), zap.Error(err))
		return fmt.Errorf("failed to delete product: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to delete product", zap.String("id", productID.String()))

	return nil
}
