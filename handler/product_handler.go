package handler

import (
	"fmt"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	IProductHandler interface {
		CreateProduct(ctx *gin.Context)
		GetAllProducts(ctx *gin.Context)
		GetProductStats(ctx *gin.Context)
		GetProductByID(ctx *gin.Context)
		GetProductBySKU(ctx *gin.Context)
		GetProductsByCategory(ctx *gin.Context)
		UpdateProduct(ctx *gin.Context)
		UpdateStock(ctx *gin.Context)
		DeleteProduct(ctx *gin.Context)
	}

	productHandler struct {
		productService service.IProductService
		logger         *zap.Logger
	}
)

func NewProductHandler(productService service.IProductService, logger *zap.Logger) *productHandler {
	return &productHandler{
		productService: productService,
		logger:         logger,
	}
}

func (ph *productHandler) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_CREATE_PRODUCT), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.CreateProduct(ctx, req)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_CREATE_PRODUCT), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_CREATE_PRODUCT), result)
	ctx.JSON(http.StatusCreated, res)
}

func (ph *productHandler) GetAllProducts(ctx *gin.Context) {
	var paginationReq dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products", dto.FAILED_GET_ALL_PRODUCTS), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PerPage <= 0 {
		paginationReq.PerPage = 10
	}

	result, err := ph.productService.GetAllProducts(ctx, paginationReq)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products", dto.FAILED_GET_ALL_PRODUCTS), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s products", dto.SUCCESS_GET_ALL_PRODUCTS), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductStats(ctx *gin.Context) {
	result, err := ph.productService.GetProductStats(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_GET_ALL_PRODUCTS, err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("success get product stats", result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductByID(ctx *gin.Context) {
	productID := ctx.Param("id")

	if productID == "" {
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCT_DETAIL, dto.MESSAGE_FAILED_INVALID_UUID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.productService.GetProductByID(ctx, productID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_PRODUCT_DETAIL), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_GET_PRODUCT_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductBySKU(ctx *gin.Context) {
	sku := ctx.Query("sku")

	if sku == "" {
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCT_DETAIL, "sku query param is required", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.productService.GetProductBySKU(ctx, sku)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_PRODUCT_DETAIL), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_GET_PRODUCT_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductsByCategory(ctx *gin.Context) {
	categoryID := ctx.Param("categoryId")

	if categoryID == "" {
		res := response.BuildResponseFailed(dto.FAILED_GET_PRODUCTS_BY_CATEGORY, dto.MESSAGE_FAILED_INVALID_UUID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var paginationReq dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products", dto.FAILED_GET_PRODUCTS_BY_CATEGORY), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PerPage <= 0 {
		paginationReq.PerPage = 10
	}

	result, err := ph.productService.GetProductsByCategoryID(ctx, categoryID, paginationReq)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products by category", dto.FAILED_GET_PRODUCTS_BY_CATEGORY), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s products by category", dto.SUCCESS_GET_PRODUCTS_BY_CATEGORY), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	if productID == "" {
		res := response.BuildResponseFailed(dto.FAILED_UPDATE_PRODUCT, dto.MESSAGE_FAILED_INVALID_UUID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE_PRODUCT), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.UpdateProduct(ctx, productID, req)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE_PRODUCT), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_UPDATE_PRODUCT), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) UpdateStock(ctx *gin.Context) {
	productID := ctx.Param("id")

	if productID == "" {
		res := response.BuildResponseFailed(dto.FAILED_UPDATE_STOCK, dto.MESSAGE_FAILED_INVALID_UUID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var req dto.UpdateStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s stock", dto.FAILED_UPDATE_STOCK), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if err := ph.productService.UpdateStock(ctx, productID, req); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE_PRODUCT), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_UPDATE_STOCK), nil)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	if productID == "" {
		res := response.BuildResponseFailed(dto.FAILED_DELETE_PRODUCT, dto.MESSAGE_FAILED_INVALID_UUID, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := ph.productService.DeleteProduct(ctx, productID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_DELETE_PRODUCT), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_DELETE_PRODUCT), nil)
	ctx.JSON(http.StatusOK, res)
}
