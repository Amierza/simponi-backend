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
	IUserHandler interface {
		GetProfile(ctx *gin.Context)
		GetLogs(ctx *gin.Context)
		GetLogsByStoreID(ctx *gin.Context)
		GetLogsByDateRange(ctx *gin.Context)

		CreateLog(ctx *gin.Context)
	}

	userHandler struct {
		userService service.IUserService
		logger      *zap.Logger
	}
)

func NewUserHandler(userService service.IUserService, logger *zap.Logger) *userHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}

func (uh *userHandler) GetProfile(ctx *gin.Context) {
	result, err := uh.userService.GetProfile(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_PROFILE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_GET_PROFILE), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) GetLogs(ctx *gin.Context) {
	result, err := uh.userService.GetLogs(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs", dto.FAILED_GET_LOGS), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s logs", dto.SUCCESS_GET_LOGS), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) GetLogsByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("storeID")
	result, err := uh.userService.GetLogsByStoreID(ctx, storeID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs by store ID", dto.FAILED_GET_LOGS_BY_STORE_ID), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s logs by store ID", dto.SUCCESS_GET_LOGS_BY_STORE_ID), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) GetLogsByDateRange(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	result, err := uh.userService.GetLogsByDateRange(ctx, startDate, endDate)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs by date range", dto.FAILED_GET_LOGS_BY_DATE_RANGE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s logs by date range", dto.SUCCESS_GET_LOGS_BY_DATE_RANGE), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) CreateLog(ctx *gin.Context) {
	var req dto.LoggingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s log", dto.FAILED_CREATE_LOG), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.CreateLog(ctx, req)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s log", dto.FAILED_CREATE_LOG), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s log", dto.SUCCESS_CREATE_LOG), result)
	ctx.JSON(http.StatusOK, res)
}
