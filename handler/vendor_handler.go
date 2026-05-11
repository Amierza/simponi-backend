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
		CreateVendorByStoreID(ctx *gin.Context)
		GetVendorsByStoreID(ctx *gin.Context)
		GetVendorByStoreIDAndVendorID(ctx *gin.Context)
		UpdateVendorByStoreIDAndVendorID(ctx *gin.Context)
		DeleteVendorByStoreIDAndVendorID(ctx *gin.Context)
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

// CreateVendor godoc
//
//	@Summary		Create vendor
//	@Description	Create a new vendor inside a store (Requires permission: CreateVendor)
//	@Tags			Vendors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string						true	"Store ID (UUID)"
//	@Param			payload		body		dto.CreateVendorRequest		true	"Create vendor request"
//	@Success		201			{object}	dto.VendorResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse			"Invalid input / UUID"
//	@Failure		401			{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse			"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse			"Store not found"
//	@Failure		500			{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/stores/{store_id}/vendors [post]
func (vh *vendorHandler) CreateVendorByStoreID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		vh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_CREATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.CreateVendorRequest
	payload.StoreID = &storeID
	if err := ctx.ShouldBind(&payload); err != nil {
		vh.logger.Error("invalid create vendor request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := vh.vendorService.CreateVendorByStoreID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

// GetVendors godoc
//
//	@Summary		Get vendors
//	@Description	Get paginated vendors in a store (Requires permission: GetVendors)
//	@Tags			Vendors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string						true	"Store ID (UUID)"
//	@Param			page		query		int							false	"Page number"
//	@Param			limit		query		int							false	"Items per page"
//	@Success		200			{object}	dto.VendorsResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse			"Invalid UUID"
//	@Failure		401			{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse			"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse			"Store not found"
//	@Failure		500			{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/stores/{store_id}/vendors [get]
func (vh *vendorHandler) GetVendorsByStoreID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		vh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_GET_ALL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		vh.logger.Error("invalid get vendors query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendors", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := vh.vendorService.GetVendorsByStoreID(ctx, &storeID, &payload)
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

// GetVendorByID godoc
//
//	@Summary		Get vendor by ID
//	@Description	Get detail vendor in a store by vendor ID (Requires permission: GetVendorByID)
//	@Tags			Vendors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string						true	"Store ID (UUID)"
//	@Param			vendor_id	path		string						true	"Vendor ID (UUID)"
//	@Success		200			{object}	dto.VendorResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse			"Invalid UUID"
//	@Failure		401			{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse			"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse			"Vendor not found"
//	@Failure		500			{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/stores/{store_id}/vendors/{vendor_id} [get]
func (vh *vendorHandler) GetVendorByStoreIDAndVendorID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		vh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	vendorIDStr := ctx.Param("vendor_id")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		vh.logger.Error("invalid vendor ID", zap.String("id", vendorIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := vh.vendorService.GetVendorByStoreIDAndVendorID(ctx, &storeID, &vendorID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateVendor godoc
//
//	@Summary		Update vendor
//	@Description	Update vendor data in a store (Requires permission: UpdateVendor)
//	@Tags			Vendors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string						true	"Store ID (UUID)"
//	@Param			vendor_id	path		string						true	"Vendor ID (UUID)"
//	@Param			payload		body		dto.UpdateVendorRequest		true	"Update vendor request"
//	@Success		200			{object}	dto.VendorResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse			"Invalid UUID / Request Body"
//	@Failure		401			{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse			"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse			"Vendor not found"
//	@Failure		500			{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/stores/{store_id}/vendors/{vendor_id} [put]
func (vh *vendorHandler) UpdateVendorByStoreIDAndVendorID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		vh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	vendorIDStr := ctx.Param("vendor_id")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		vh.logger.Error("invalid vendor ID", zap.String("id", vendorIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateVendorRequest
	payload.StoreID = &storeID
	payload.ID = vendorID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		vh.logger.Error("invalid update vendor request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := vh.vendorService.UpdateVendorByStoreIDAndVendorID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

// DeleteVendorByID godoc
//
//	@Summary		Delete vendor
//	@Description	Delete vendor from a store by vendor ID (Requires permission: DeleteVendorByID)
//	@Tags			Vendors
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id	path		string							true	"Store ID (UUID)"
//	@Param			vendor_id	path		string							true	"Vendor ID (UUID)"
//	@Success		200			{object}	dto.EmptySuccessResponseWrapper	"Success"
//	@Failure		400			{object}	dto.ErrorResponse				"Invalid UUID"
//	@Failure		401			{object}	dto.ErrorResponse				"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse				"Forbidden"
//	@Failure		404			{object}	dto.ErrorResponse				"Vendor not found"
//	@Failure		500			{object}	dto.ErrorResponse				"Internal Server Error"
//	@Router			/stores/{store_id}/vendors/{vendor_id} [delete]
func (vh *vendorHandler) DeleteVendorByStoreIDAndVendorID(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		vh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	vendorIDStr := ctx.Param("vendor_id")
	vendorID, err := uuid.Parse(vendorIDStr)
	if err != nil {
		vh.logger.Error("invalid vendor ID", zap.String("id", vendorIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := vh.vendorService.DeleteVendorByStoreIDAndVendorID(ctx, &storeID, &vendorID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s vendor", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s vendor", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
