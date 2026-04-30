package handler

import (
	"fmt"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IProductHandler interface {
		CreateProduct(ctx *gin.Context)
		GetProducts(ctx *gin.Context)
		GetProductCategory(ctx *gin.Context)
		GetProductStats(ctx *gin.Context)
		GetProductByID(ctx *gin.Context)
		UpdateProduct(ctx *gin.Context)
		UpdateStock(ctx *gin.Context)
		DeleteProductByID(ctx *gin.Context)
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
	var payload dto.CreateProductRequest
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ph.logger.Error("invalid create product request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.CreateProduct(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

func (ph *productHandler) GetProducts(ctx *gin.Context) {
	var payload dto.ProductPaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		ph.logger.Error("invalid get products query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.GetProducts(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s products", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductCategory(ctx *gin.Context) {
	result, err := ph.productService.GetProductCategory(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product category", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product category", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductStats(ctx *gin.Context) {
	result, err := ph.productService.GetProductStats(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product stats", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product stats", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductByID(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)

	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.productService.GetProductByID(ctx, &productID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) UpdateProduct(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ph.logger.Error("invalid update product request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.UpdateProduct(ctx, &productID, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) UpdateStock(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateStockRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ph.logger.Error("invalid update stock request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s stock", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if err := ph.productService.UpdateStock(ctx, &productID, &payload); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s stock", dto.SUCCESS_UPDATE), nil)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) DeleteProductByID(ctx *gin.Context) {
	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := ph.productService.DeleteProductByID(ctx, &productID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
