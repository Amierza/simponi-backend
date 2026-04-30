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
	IStoreUserHandler interface {
		CreateStoreUsers(ctx *gin.Context)
		GetStoreUsers(ctx *gin.Context)
		GetStoreUserByStoreIDAndUserID(ctx *gin.Context)
		DeleteStoreUserByStoreIDAndUserID(ctx *gin.Context)
	}

	storeUserHandler struct {
		storeUserService service.IStoreUserService
		logger           *zap.Logger
	}
)

func NewStoreUserHandler(storeUserService service.IStoreUserService, logger *zap.Logger) *storeUserHandler {
	return &storeUserHandler{
		storeUserService: storeUserService,
		logger:           logger,
	}
}

func (suh *storeUserHandler) CreateStoreUsers(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		suh.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store users", dto.FAILED_CREATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.CreateStoreUsersRequest
	payload.StoreID = &storeID
	if err := ctx.ShouldBind(&payload); err != nil {
		suh.logger.Error("invalid create store user request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store users", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	err = suh.storeUserService.CreateStoreUsers(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store users", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store users", dto.SUCCESS_CREATE), nil)
	ctx.JSON(http.StatusCreated, res)
}

func (suh *storeUserHandler) GetStoreUsers(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		suh.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store users", dto.FAILED_GET_ALL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		suh.logger.Error("invalid get store users query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store users", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := suh.storeUserService.GetStoreUsers(ctx, &payload, &storeID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store users", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s store users", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (suh *storeUserHandler) GetStoreUserByStoreIDAndUserID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		suh.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store user", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		suh.logger.Error("invalid user ID", zap.String("user_id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store user", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := suh.storeUserService.GetStoreUserByStoreIDAndUserID(ctx, &storeID, &userID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store user", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store user", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (suh *storeUserHandler) DeleteStoreUserByStoreIDAndUserID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		suh.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store user", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		suh.logger.Error("invalid user ID", zap.String("user_id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s store user", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := suh.storeUserService.DeleteStoreUserByStoreIDAndUserID(ctx, &storeID, &userID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s store user", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s store user", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
