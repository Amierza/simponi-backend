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
	storeID := ctx.Query("store_id")
	startDate := ctx.Query("start_date")
	if startDate == "" {
		startDate = ctx.Query("startdate")
	}
	endDate := ctx.Query("end_date")
	if endDate == "" {
		endDate = ctx.Query("enddate")
	}

	var paginationReq response.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s", dto.FAILED_GET_LOGS), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if paginationReq.Page <= 0 {
		paginationReq.Page = 1
	}
	if paginationReq.PerPage <= 0 {
		paginationReq.PerPage = 10
	}

	result, err := lh.logService.GetLogs(ctx, storeID, startDate, endDate, paginationReq)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s", dto.FAILED_GET_LOGS), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s", dto.SUCCESS_GET_LOGS),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (lh *logHandler) CreateLog(ctx *gin.Context) {
	var req dto.LogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		status := http.StatusBadRequest
		res := response.BuildResponseFailed(fmt.Sprintf("%s", dto.FAILED_CREATE_LOG), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := lh.logService.CreateLog(ctx, req)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s", dto.FAILED_CREATE_LOG), err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s", dto.SUCCESS_CREATE_LOG), result)
	ctx.JSON(http.StatusOK, res)
}
