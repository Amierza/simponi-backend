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
	IExternalProductHandler interface {
		CreateExternalProduct(ctx *gin.Context)
		GetExternalProducts(ctx *gin.Context)
		GetExternalProductByStoreIDAndExprodID(ctx *gin.Context)
		GetExternalProductsByStoreIDAndStorePlatformID(ctx *gin.Context)
		UpdateExternalProductByStoreIDAndExprodID(ctx *gin.Context)
		DeleteExternalProductByStoreIDAndExprodID(ctx *gin.Context)
	}

	externalProductHandler struct {
		externalProductService service.IExternalProductService
		logger                 *zap.Logger
	}
)

func NewExternalProductHandler(externalProductService service.IExternalProductService, logger *zap.Logger) *externalProductHandler {
	return &externalProductHandler{
		externalProductService: externalProductService,
		logger:                 logger,
	}
}

// CreateExternalProduct godoc
//
//	@Summary		Create external product
//	@Description	Create a new external product for a store
//	@Tags			External Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			store_id	path		string								true	"Store ID"
//	@Param			request		body		dto.CreateExternalProductRequest	true	"Request Body"
//	@Success		201			{object}	dto.ExternalProductCreateSuccessResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Router			/stores/{store_id}/external-products [post]
func (eph *externalProductHandler) CreateExternalProduct(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		eph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.CreateExternalProductRequest
	payload.StoreID = &storeID
	if err := ctx.ShouldBind(&payload); err != nil {
		eph.logger.Error("invalid create external product request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := eph.externalProductService.CreateExternalProduct(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

// GetExternalProducts godoc
//
//	@Summary		Get external products
//	@Description	Get list of external products by store
//	@Tags			External Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			store_id	path		string	true	"Store ID"
//	@Param			page		query		int		false	"Page"
//	@Param			limit		query		int		false	"Limit"
//	@Success		200			{object}	dto.ExternalProductsSuccessResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Router			/stores/{store_id}/external-products [get]
func (eph *externalProductHandler) GetExternalProducts(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		eph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		eph.logger.Error("invalid get external products query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProducts(ctx, &payload, &storeID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s external products", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

// GetExternalProductByStoreIDAndExprodID godoc
//
//	@Summary		Get external product detail
//	@Description	Get detail of external product by ID
//	@Tags			External Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			store_id			path		string	true	"Store ID"
//	@Param			external_product_id	path		string	true	"External Product ID"
//	@Success		200					{object}	dto.ExternalProductSuccessResponse
//	@Failure		400					{object}	dto.ErrorResponse
//	@Router			/stores/{store_id}/external-products/{external_product_id} [get]
func (eph *externalProductHandler) GetExternalProductByStoreIDAndExprodID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		eph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	externalProductIDStr := ctx.Param("external_product_id")
	externalProductID, err := uuid.Parse(externalProductIDStr)
	if err != nil {
		eph.logger.Error("invalid external product ID", zap.String("id", externalProductIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProductByStoreIDAndExprodID(ctx, &storeID, &externalProductID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

// GetExternalProductsByStoreIDAndStorePlatformID godoc
//
//	@Summary		Get external products by platform
//	@Description	Get external products filtered by store platform
//	@Tags			External Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			store_id			path		string	true	"Store ID"
//	@Param			store_platform_id	path		string	true	"Store Platform ID"
//	@Success		200					{object}	dto.ExternalProductsSuccessResponse
//	@Failure		400					{object}	dto.ErrorResponse
//	@Router			/stores/{store_id}/external-products/store-platform/{store_platform_id} [get]
func (eph *externalProductHandler) GetExternalProductsByStoreIDAndStorePlatformID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		eph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	storePlatformIDStr := ctx.Param("store_platform_id")
	storePlatformID, err := uuid.Parse(storePlatformIDStr)
	if err != nil {
		eph.logger.Error("invalid store platform ID", zap.String("store_platform_id", storePlatformIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProductsByStoreIDAndStorePlatformID(ctx, &storeID, &storePlatformID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external products", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateExternalProductByStoreIDAndExprodID godoc
//
//	@Summary		Update external product
//	@Description	Update external product price
//	@Tags			External Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			store_id			path		string								true	"Store ID"
//	@Param			external_product_id	path		string								true	"External Product ID"
//	@Param			request				body		dto.UpdateExternalProductRequest	true	"Request Body"
//	@Success		200					{object}	dto.ExternalProductUpdateSuccessResponse
//	@Failure		400					{object}	dto.ErrorResponse
//	@Router			/stores/{store_id}/external-products/{external_product_id} [put]
func (eph *externalProductHandler) UpdateExternalProductByStoreIDAndExprodID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		eph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	externalProductIDString := ctx.Param("external_product_id")
	externalProductID, err := uuid.Parse(externalProductIDString)
	if err != nil {
		eph.logger.Error("invalid external product ID", zap.String("id", externalProductIDString), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateExternalProductRequest
	payload.StoreID = &storeID
	payload.ID = externalProductID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		eph.logger.Error("invalid update external product request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := eph.externalProductService.UpdateExternalProductByStoreIDAndExprodID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

// DeleteExternalProductByStoreIDAndExprodID godoc
//
//	@Summary		Delete external product
//	@Description	Delete external product by ID
//	@Tags			External Products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			store_id			path		string	true	"Store ID"
//	@Param			external_product_id	path		string	true	"External Product ID"
//	@Success		200					{object}	dto.ExternalProductDeleteSuccessResponse
//	@Failure		400					{object}	dto.ErrorResponse
//	@Router			/stores/{store_id}/external-products/{external_product_id} [delete]
func (eph *externalProductHandler) DeleteExternalProductByStoreIDAndExprodID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		eph.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	externalProductIDStr := ctx.Param("external_product_id")
	externalProductID, err := uuid.Parse(externalProductIDStr)
	if err != nil {
		eph.logger.Error("invalid external product ID", zap.String("id", externalProductIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := eph.externalProductService.DeleteExternalProductByStoreIDAndExprodID(ctx, &storeID, &externalProductID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
