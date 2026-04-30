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
	IStoreHandler interface {
		CreateStore(ctx *gin.Context)
		GetStores(ctx *gin.Context)
		GetStoreByStoreID(ctx *gin.Context)
		UpdateStoreByStoreID(ctx *gin.Context)
		DeleteStoreByStoreID(ctx *gin.Context)
	}

	storeHandler struct {
		storeService service.IStoreService
		logger       *zap.Logger
	}
)

func NewStoreHandler(storeService service.IStoreService, logger *zap.Logger) *storeHandler {
	return &storeHandler{
		storeService: storeService,
		logger:       logger,
	}
}

func (sh *storeHandler) CreateStore(ctx *gin.Context) {
	var payload dto.CreateStoreRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		sh.logger.Error("invalid create store request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := sh.storeService.CreateStore(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

func (sh *storeHandler) GetStores(ctx *gin.Context) {
	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		sh.logger.Error("invalid get stores query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s stores", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := sh.storeService.GetStores(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s stores", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s stores", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (sh *storeHandler) GetStoreByStoreID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		sh.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sh.storeService.GetStoreByStoreID(ctx, &storeID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (sh *storeHandler) UpdateStoreByStoreID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		sh.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateStoreRequest
	payload.ID = storeID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		sh.logger.Error("invalid update store request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := sh.storeService.UpdateStoreByStoreID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (sh *storeHandler) DeleteStoreByStoreID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		sh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := sh.storeService.DeleteStoreByStoreID(ctx, &storeID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
