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
		GetExternalProductByID(ctx *gin.Context)
		GetExternalProductsByProductID(ctx *gin.Context)
		GetExternalProductsByStorePlatformID(ctx *gin.Context)
		UpdateExternalProduct(ctx *gin.Context)
		DeleteExternalProductByID(ctx *gin.Context)
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

func (eph *externalProductHandler) CreateExternalProduct(ctx *gin.Context) {
	var payload dto.CreateExternalProductRequest
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

func (eph *externalProductHandler) GetExternalProducts(ctx *gin.Context) {
	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		eph.logger.Error("invalid get external products query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProducts(ctx, &payload)
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

func (eph *externalProductHandler) GetExternalProductByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		eph.logger.Error("invalid external product ID", zap.String("id", idStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProductByID(ctx, &id)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (eph *externalProductHandler) GetExternalProductsByProductID(ctx *gin.Context) {
	productIDStr := ctx.Param("productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		eph.logger.Error("invalid product ID", zap.String("productId", productIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProductsByProductID(ctx, &productID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external products", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

func (eph *externalProductHandler) GetExternalProductsByStorePlatformID(ctx *gin.Context) {
	storePlatformIDStr := ctx.Param("storePlatformId")
	storePlatformID, err := uuid.Parse(storePlatformIDStr)
	if err != nil {
		eph.logger.Error("invalid store platform ID", zap.String("storePlatformId", storePlatformIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := eph.externalProductService.GetExternalProductsByStorePlatformID(ctx, &storePlatformID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external products", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external products", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

func (eph *externalProductHandler) UpdateExternalProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		eph.logger.Error("invalid external product ID", zap.String("id", idStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateExternalProductRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		eph.logger.Error("invalid update external product request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := eph.externalProductService.UpdateExternalProduct(ctx, &id, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (eph *externalProductHandler) DeleteExternalProductByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		eph.logger.Error("invalid external product ID", zap.String("id", idStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := eph.externalProductService.DeleteExternalProductByID(ctx, &id); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s external product", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s external product", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}