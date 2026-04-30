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
		GetProductStats(ctx *gin.Context)
		GetProductByStoreIDAndProductID(ctx *gin.Context)
		UpdateProductByStoreIDAndProductID(ctx *gin.Context)
		UpdateStockByStoreIDAndProductID(ctx *gin.Context)
		DeleteProductByStoreIDAndProductID(ctx *gin.Context)
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
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.CreateProductRequest
	payload.StoreID = &storeID
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
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		ph.logger.Error("invalid get products query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.GetProducts(ctx, &payload, &storeID)
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

func (ph *productHandler) GetProductStats(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.productService.GetProductStats(ctx, &storeID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product stats", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product stats", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) GetProductByStoreIDAndProductID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.productService.GetProductByStoreIDAndProductID(ctx, &storeID, &productID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) UpdateProductByStoreIDAndProductID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateProductRequest
	payload.StoreID = &storeID
	payload.ID = productID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ph.logger.Error("invalid update product request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.productService.UpdateProductByStoreIDAndProductID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) UpdateStockByStoreIDAndProductID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateStockRequest
	payload.StoreID = &storeID
	payload.ID = productID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ph.logger.Error("invalid update stock request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s stock", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if err := ph.productService.UpdateStockByStoreIDAndProductID(ctx, &payload); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s stock", dto.SUCCESS_UPDATE), nil)
	ctx.JSON(http.StatusOK, res)
}

func (ph *productHandler) DeleteProductByStoreIDAndProductID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		ph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		ph.logger.Error("invalid product ID", zap.String("id", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := ph.productService.DeleteProductByStoreIDAndProductID(ctx, &storeID, &productID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s product", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
