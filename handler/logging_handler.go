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
	ILogHandler interface {
		CreateLog(ctx *gin.Context)
		GetLogs(ctx *gin.Context)
		GetLogsByStoreID(ctx *gin.Context)
		GetLogsByDateRange(ctx *gin.Context)
	}

	logHandler struct {
		logService service.ILogService
		logger     *zap.Logger
	}
)

func NewLogHandler(logService service.ILogService, logger *zap.Logger) *logHandler {
	return &logHandler{
		logService: logService,
		logger:     logger,
	}
}

func (lh *logHandler) GetLogs(ctx *gin.Context) {
	var paginationReq response.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs", dto.FAILED_GET_LOGS), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PerPage <= 0 {
		paginationReq.PerPage = 10
	}

	result, err := lh.logService.GetLogs(ctx, paginationReq)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs", dto.FAILED_GET_LOGS), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s logs", dto.SUCCESS_GET_LOGS), result)
	ctx.JSON(http.StatusOK, res)
}

func (lh *logHandler) GetLogsByStoreID(ctx *gin.Context) {
	storeID := ctx.Query("storeID")

	var paginationReq response.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs", dto.FAILED_GET_LOGS), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PerPage <= 0 {
		paginationReq.PerPage = 10
	}

	result, err := lh.logService.GetLogsByStoreID(ctx, storeID, paginationReq)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs by store ID", dto.FAILED_GET_LOGS_BY_STORE_ID), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s logs by store ID", dto.SUCCESS_GET_LOGS_BY_STORE_ID), result)
	ctx.JSON(http.StatusOK, res)
}

func (lh *logHandler) GetLogsByDateRange(ctx *gin.Context) {
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	var paginationReq response.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs", dto.FAILED_GET_LOGS), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PerPage <= 0 {
		paginationReq.PerPage = 10
	}

	result, err := lh.logService.GetLogsByDateRange(ctx, startDate, endDate, paginationReq)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s logs by date range", dto.FAILED_GET_LOGS_BY_DATE_RANGE), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s logs by date range", dto.SUCCESS_GET_LOGS_BY_DATE_RANGE), result)
	ctx.JSON(http.StatusOK, res)
}

func (lh *logHandler) CreateLog(ctx *gin.Context) {
	var req dto.LogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s log", dto.FAILED_CREATE_LOG), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := lh.logService.CreateLog(ctx, req)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s log", dto.FAILED_CREATE_LOG), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s log", dto.SUCCESS_CREATE_LOG), result)
	ctx.JSON(http.StatusOK, res)
}
