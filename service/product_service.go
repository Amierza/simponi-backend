package service

import (
	"context"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IProductService interface {
		CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error)
		GetAllProducts(ctx context.Context, req dto.PaginationRequest) (dto.ProductPaginationResponse, error)
		GetProductByID(ctx context.Context, productID string) (dto.ProductResponse, error)
		GetProductBySKU(ctx context.Context, sku string) (dto.ProductResponse, error)
		GetProductsByCategoryID(ctx context.Context, categoryID string, req dto.PaginationRequest) (dto.ProductResponse, error)
		UpdateProduct(ctx context.Context, productID string, req dto.UpdateProductRequest) (dto.ProductResponse, error)
		DeleteProduct(ctx context.Context, productID string) error

		UpdateStock(ctx context.Context, productID string, req dto.UpdateStockRequest) error
	}

	productService struct {
		productRepo	repository.IProductRepository
		logger		*zap.Logger
		jwtService	jwt.IJWT
	}
)

func NewProductService(productRepo repository.IProductRepository, logger *zap.Logger, jwtService jwt.IJWT) *productService {
	return &productService{
		productRepo:	productRepo,
		logger:			logger,
		jwtService: 	jwtService,
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

func mapToProductListResponse(p entity.Product) dto.ProductListResponse {
	var category *dto.ProductCategoryResponse
	if p.CategoryID != nil {
		category = &dto.ProductCategoryResponse{
			ID:		p.Category.ID,
			Name:	p.Category.Name,
		}
	}

	var images []dto.ProductImageResponse
	for _, img := range p.Images {
		images = append(images, dto.ProductImageResponse{
			ID: 		img.ID,
			ImageURL: 	img.ImageURL,
		})
	}

	var externalProducts []dto.ExternalProductResponse
	for _, ep := range p.ExternalProducts {
		externalProducts = append(externalProducts, dto.ExternalProductResponse{
			ID: 				ep.ID,
			ProductID:			ep.ProductID,
			StoreID:			ep.StoreID,
			ExternalProductID:	ep.ExternalProductID,
			ExternalModelID:	ep.ExternalModelID,
			Price:				ep.Price,
		})
	}

	return dto.ProductListResponse{
		ID:					p.ID,
		Name:				p.Name,
		SKU:				p.SKU,
		Stock: 				p.Stock,
		Category:			category,
		Images:				images,
		ExternalProducts:	externalProducts,
		Status:				getProductStatus(p),
		CreatedAt:			p.CreatedAt,
	}
}

func mapToProductResponse(p entity.Product) dto.ProductResponse {
	var category *dto.ProductCategoryResponse
	if p.CategoryID != nil {
		category = &dto.ProductCategoryResponse{
			ID:		p.Category.ID,
			Name:	p.Category.Name,
		}
	}

	var images []dto.ProductImageResponse
	for _, img := range p.Images {
		images = append(images, dto.ProductImageResponse{
			ID: 		img.ID,
			ImageURL: 	img.ImageURL,
		})
	}

	var externalProducts []dto.ExternalProductResponse
	for _, ep := range p.ExternalProducts {
		externalProducts = append(externalProducts, dto.ExternalProductResponse{
			ID: 				ep.ID,
			ProductID:			ep.ProductID,
			StoreID:			ep.StoreID,
			ExternalProductID:	ep.ExternalProductID,
			ExternalModelID:	ep.ExternalModelID,
			Price:				ep.Price,
		})
	}

	return dto.ProductResponse{
		ID:					p.ID,
		Name:				p.Name,
		Description:		p.Description,
		SKU:				p.SKU,
		Stock: 				p.Stock,
		Category:			category,
		Images:				images,
		ExternalProducts:	externalProducts,
		CreatedAt:			p.CreatedAt,
		UpdatedAt: 			p.UpdatedAt,
	}
}

func (ps *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	existing, _ := ps.productRepo.GetProductBySKU(ctx, nil, req.SKU)
	if existing != nil {
		ps.logger.Warn("Product SKU already exists", zap.String("sku", req.SKU))
		return dto.ProductResponse{}, dto.ErrProductSKUAlreadyExists
	}

	product, err := ps.productRepo.CreateProduct(ctx, nil, &entity.Product{
		ID:				uuid.New(),
		Name:			req.Name,
		Description:	req.Description,
		SKU:			req.SKU,
		Stock:			req.Stock,
		CategoryID: 	req.CategoryID,
	})

	if err != nil {
		ps.logger.Error("failed to create product", zap.Error(err))
		return dto.ProductResponse{}, dto.ErrCreateProduct 
	}

	return mapToProductResponse(*product), nil
}

func (ps *productService) GetAllProducts(ctx context.Context, req dto.PaginationRequest) (dto.ProductPaginationResponse, error) {
	products, err := ps.productRepo.GetProducts(ctx, nil)
	
	if err != nil {
		ps.logger.Error("failed to get all products", zap.Error(err))
		return dto.ProductPaginationResponse{}, dto.ErrGetAllProducts
	}

	if len(products) == 0 {
		ps.logger.Warn("products not found")
		return dto.ProductPaginationResponse{}, dto.ErrNotFound
	}

	var productList []dto.ProductListResponse
	for _, p := range products {
		productList = append(productList, mapToProductListResponse(p))
	}

	count := int64(len(productList))
	perPage := int64(req.PerPage)
	maxPage := (count + perPage - 1) / perPage

	return dto.ProductPaginationResponse{
		Data: productList,
		Pagination: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (ps *productService) GetProductByID(ctx context.Context, productID string) (dto.ProductResponse, error) {
	product, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	
	if err != nil {
		ps.logger.Error("failed to get product by ID", zap.String("productID", productID), zap.Error(err))
		return dto.ProductResponse{}, dto.ErrGetProductByID
	}

	if product == nil {
		ps.logger.Error("product not found", zap.String("productID", productID))
		return dto.ProductResponse{}, dto.ErrGetProductByID
	}

	return mapToProductResponse(*product), nil
}

func (ps *productService) GetProductBySKU(ctx context.Context, sku string) (dto.ProductResponse, error) {
	product, err := ps.productRepo.GetProductByID(ctx, nil, sku)
	
	if err != nil {
		ps.logger.Error("failed to get product by SKU", zap.String("sku", sku), zap.Error(err))
		return dto.ProductResponse{}, dto.ErrGetProductBySKU
	}

	if product == nil {
		ps.logger.Error("product not found", zap.String("sku", sku))
		return dto.ProductResponse{}, dto.ErrGetProductBySKU
	}

	return mapToProductResponse(*product), nil
}

func (ps *productService) GetProductsByCategory(ctx context.Context, categoryID string, req dto.PaginationRequest) (dto.ProductPaginationResponse, error) {
	products, err := ps.productRepo.GetProductsByCategoryID(ctx, nil, categoryID)
	
	if err != nil {
		ps.logger.Error("failed to get products by category", zap.String("categoryID", categoryID), zap.Error(err))
		return dto.ProductPaginationResponse{}, dto.ErrGetProductsByCategory
	}

	if len(products) == 0 {
		ps.logger.Warn("no products found for category", zap.String("catergoryID", categoryID))
		return dto.ProductPaginationResponse{}, dto.ErrGetProductsByCategory
	}

	var productList []dto.ProductListResponse
	for _, p := range products {
		productList = append(productList, mapToProductListResponse(p))
	}

	count := int64(len(productList))
	perPage := int64(req.PerPage)
	maxPage := (count + perPage - 1) / perPage

	return dto.ProductPaginationResponse{
		Data: productList,
		Pagination: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (ps *productService) UpdateProduct(ctx context.Context, productID string, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	product, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		ps.logger.Warn("Product not found for update", zap.String("productID", productID))
		return dto.ProductResponse{}, dto.ErrNotFound
	}

	if req.SKU != "" && req.SKU != product.SKU {
		existing, _ := ps.productRepo.GetProductBySKU(ctx, nil, req.SKU)
		if existing != nil {
			ps.logger.Warn("SKU already taken", zap.String("sku", req.SKU))
			return dto.ProductResponse{}, dto.ErrProductSKUAlreadyExists
		}
		product.SKU = req.SKU
	}

	if req.Name != "" {
		product.Name = req.Name
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}

	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	updated, err := ps.productRepo.UpdateProduct(ctx, nil, product)

	if err != nil {
		ps.logger.Error("failed to update product", zap.String("productID", productID), zap.Error(err))
		return dto.ProductResponse{}, dto.ErrUpdateProduct
	}

	return mapToProductResponse(*updated), nil
}

func (ps *productService) UpdateStock(ctx context.Context, productID string, req dto.UpdateStockRequest) error {
	_, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		ps.logger.Warn("Product not found for stock update", zap.String("productID", productID), zap.Error(err))
		return dto.ErrNotFound
	}

	if err := ps.productRepo.UpdateStock(ctx, nil, productID, req.Change); err != nil {
		ps.logger.Error("failed to update stock", zap.String("productID", productID), zap.Error(err))
		return dto.ErrUpdateProduct
	}

	return nil
}

func (ps *productService) DeleteStock(ctx context.Context, productID string) error {
	_, err := ps.productRepo.GetProductByID(ctx, nil, productID)
	if err != nil {
		ps.logger.Warn("Product not found for delete", zap.String("productID", productID), zap.Error(err))
		return dto.ErrNotFound
	}

	if err := ps.productRepo.DeleteProduct(ctx, nil, productID); err != nil {
		ps.logger.Error("failed to delete stock", zap.String("productID", productID), zap.Error(err))
		return dto.ErrDeleteProduct
	}

	return nil
}