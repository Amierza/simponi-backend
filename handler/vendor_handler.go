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
	IVendorHandler interface {
		CreateVendor(ctx *gin.Context)
		GetVendors(ctx *gin.Context)
		GetVendorByID(ctx *gin.Context)
		UpdateVendor(ctx *gin.Context)
		DeleteVendorByID(ctx *gin.Context)
	}

	vendorHandler struct {
		vendorService service.IVendorService
		logger        *zap.Logger
	}
)

func NewVendorHandler(vendorService service.IVendorService, logger *zap.Logger) *vendorHandler {
	return &vendorHandler{
		vendorService: vendorService,
		logger:        logger,
	}
}

func (vh *vendorHandler) CreateVendor(ctx *gin.Context) {
	var payload dto.CreateVendorRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		vh.logger.Error("invalid create vendor request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := vh.vendorService.CreateVendor(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

func (vh *vendorHandler) GetVendors(ctx *gin.Context) {
	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		vh.logger.Error("invalid get vendors query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendors", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := vh.vendorService.GetVendors(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendors", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s vendors", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (vh *vendorHandler) GetVendorByID(ctx *gin.Context) {
	vendorIDStr := ctx.Param("id")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		vh.logger.Error("invalid vendor ID", zap.String("id", vendorIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := vh.vendorService.GetVendorByID(ctx, &vendorID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (vh *vendorHandler) UpdateVendor(ctx *gin.Context) {
	vendorIDStr := ctx.Param("id")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		vh.logger.Error("invalid vendor ID", zap.String("id", vendorIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateVendorRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		vh.logger.Error("invalid update vendor request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := vh.vendorService.UpdateVendor(ctx, &vendorID, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (vh *vendorHandler) DeleteVendorByID(ctx *gin.Context) {
	vendorIDStr := ctx.Param("id")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		vh.logger.Error("invalid vendor ID", zap.String("id", vendorIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := vh.vendorService.DeleteVendorByID(ctx, &vendorID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
