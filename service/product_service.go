package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IProductService interface {
		CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (dto.ProductResponse, error)
		GetProducts(ctx context.Context, req *dto.ProductPaginationRequest) (dto.ProductPaginationResponse, error)
		GetProductCategory(ctx context.Context) ([]dto.ProductCategoryResponse, error)
		GetProductStats(ctx context.Context) (dto.ProductStatsResponse, error)
		GetProductByID(ctx context.Context, productID *uuid.UUID) (dto.ProductResponse, error)
		GetProductBySKU(ctx context.Context, sku string) (dto.ProductResponse, error)
		UpdateProduct(ctx context.Context, productID *uuid.UUID, req *dto.UpdateProductRequest) (dto.ProductResponse, error)
		DeleteProductByID(ctx context.Context, productID *uuid.UUID) error

		UpdateStock(ctx context.Context, productID *uuid.UUID, req *dto.UpdateStockRequest) error
	}

	productService struct {
		productRepo repository.IProductRepository
		logger      *zap.Logger
		jwtService  jwt.IJWT
	}
)

func NewProductService(productRepo repository.IProductRepository, logger *zap.Logger, jwtService jwt.IJWT) *productService {
	return &productService{
		productRepo: productRepo,
		logger:      logger,
		jwtService:  jwtService,
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
		Category:         MapToProductCategoryResponse(p),
		Images:           MapToProductImageResponse(p),
		ExternalProducts: MapToExternalProductResponse(p),
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func (ps *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (dto.ProductResponse, error) {
	_, found, err := ps.productRepo.GetProductBySKU(ctx, nil, req.SKU)
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
	})

	if err != nil {
		ps.logger.Error("failed to create product", zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to create product: %w", dto.ErrInternal)
	}

	if err := ps.productRepo.AttachProductImageToProduct(ctx, nil, req.ImageID, &newProduct.ID); err != nil {
		ps.logger.Error("failed to attach image to product", zap.String("productID", newProduct.ID.String()), zap.String("imageID", req.ImageID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to attach image to product: %w", dto.ErrBadRequest)
	}

	newProduct, _, err = ps.productRepo.GetProductByID(ctx, nil, &newProduct.ID)
	if err != nil {
		ps.logger.Error("failed to reload product after image attach", zap.String("productID", newProduct.ID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to load product: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to create product", zap.String("id", newProduct.ID.String()))

	return MapToProductResponse(*newProduct), nil
}

func (ps *productService) GetProducts(ctx context.Context, req *dto.ProductPaginationRequest) (dto.ProductPaginationResponse, error) {
	datas, err := ps.productRepo.GetProducts(ctx, nil, req)

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

func (ps *productService) GetProductCategory(ctx context.Context) ([]dto.ProductCategoryResponse, error) {
	categories, err := ps.productRepo.GetProductCategory(ctx, nil)
	if err != nil {
		ps.logger.Error("failed to get product categories", zap.Error(err))
		return nil, fmt.Errorf("failed to get product stats: %w", dto.ErrInternal)
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


func (ps *productService) GetProductStats(ctx context.Context) (dto.ProductStatsResponse, error) {
	stats, err := ps.productRepo.GetProductStats(ctx, nil)
	if err != nil {
		ps.logger.Error("failed to get product stats", zap.Error(err))
		return dto.ProductStatsResponse{}, fmt.Errorf("failed to get product stats: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to get product stats")

	return stats, nil
}

func (ps *productService) GetProductByID(ctx context.Context, productID *uuid.UUID) (dto.ProductResponse, error) {
	product, found, err := ps.productRepo.GetProductByID(ctx, nil, productID)

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

func (ps *productService) GetProductBySKU(ctx context.Context, sku string) (dto.ProductResponse, error) {
	product, found, err := ps.productRepo.GetProductBySKU(ctx, nil, sku)

	if err != nil {
		ps.logger.Error("failed to get product by SKU", zap.String("sku", sku), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to get product by SKU: %w", dto.ErrInternal)
	}

	if !found {
		ps.logger.Warn("product not found", zap.String("sku", sku))
		return dto.ProductResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	ps.logger.Info("success to get product by SKU", zap.String("sku", sku))

	return MapToProductResponse(*product), nil
}

func (ps *productService) UpdateProduct(ctx context.Context, productID *uuid.UUID, req *dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, found, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", productID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}

	if !found {
		ps.logger.Warn("product not found", zap.String("productID", productID.String()))
		return dto.ProductResponse{}, fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	if req.SKU != "" && req.SKU != product.SKU {
		_, found, err := ps.productRepo.GetProductBySKU(ctx, nil, req.SKU)
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

	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	updatedProduct, err := ps.productRepo.UpdateProduct(ctx, nil, product)

	if err != nil {
		ps.logger.Error("failed to update product", zap.String("id", productID.String()), zap.Error(err))
		return dto.ProductResponse{}, fmt.Errorf("failed to update product: %w", dto.ErrInternal)
	}

	return MapToProductResponse(*updatedProduct), nil
}

func (ps *productService) UpdateStock(ctx context.Context, productID *uuid.UUID, req *dto.UpdateStockRequest) error {
	_, found, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", productID.String()), zap.Error(err))
		return fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}

	if !found {
		ps.logger.Warn("product not found", zap.String("productID", productID.String()))
		return fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	if err := ps.productRepo.UpdateStock(ctx, nil, productID, req.Change); err != nil {
		ps.logger.Error("failed to update stock", zap.String("productID", productID.String()), zap.Error(err))
		return dto.ErrInternal
	}

	ps.logger.Info("success to update stock", zap.String("categoryID", productID.String()), zap.Int("change", req.Change))

	return nil
}

func (ps *productService) DeleteProductByID(ctx context.Context, productID *uuid.UUID) error {
	_, found, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", productID.String()), zap.Error(err))
		return fmt.Errorf("failed to get product by ID: %w", dto.ErrInternal)
	}

	if !found {
		ps.logger.Warn("product not found", zap.String("productID", productID.String()))
		return fmt.Errorf("product not found: %w", dto.ErrNotFound)
	}

	if err := ps.productRepo.DeleteProductByID(ctx, nil, productID); err != nil {
		ps.logger.Error("failed to delete product", zap.String("productID", productID.String()), zap.Error(err))
		return fmt.Errorf("product not found: %w", dto.ErrInternal)
	}

	ps.logger.Info("success to delete product", zap.String("id", productID.String()))

	return nil
}
