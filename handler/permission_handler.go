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
	IPermissionHandler interface {
		GetPermissions(ctx *gin.Context)
	}

	permissionHandler struct {
		permissionService service.IPermissionService
		logger            *zap.Logger
	}
)

func NewPermissionHandler(permissionService service.IPermissionService, logger *zap.Logger) *permissionHandler {
	return &permissionHandler{
		permissionService: permissionService,
		logger:            logger,
	}
}

// GetPermissions godoc
//
//	@Summary		Get list of permissions
//	@Description	Get all permissions or paginated permissions (Requires permission: GetPermissions)
//	@Description	Set query param `pagination=false` to disable pagination
//	@Tags			Permissions
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			pagination	query		string									false	"Enable pagination (true/false), default: true"
//	@Param			page		query		int										false	"Page number (if pagination=true)"
//	@Param			limit		query		int										false	"Items per page (if pagination=true)"
//	@Success		200			{object}	dto.PermissionPaginationResponseWrapper	"Success (with pagination)"
//	@Success		200			{object}	dto.PermissionResponseWrapper			"Success (without pagination)"
//	@Failure		400			{object}	dto.ErrorResponse						"Bad Request"
//	@Failure		401			{object}	dto.ErrorResponse						"Unauthorized"
//	@Failure		403			{object}	dto.ErrorResponse						"Forbidden"
//	@Failure		500			{object}	dto.ErrorResponse						"Internal Server Error"
//	@Router			/permissions [get]
func (ph *permissionHandler) GetPermissions(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		result, err := ph.permissionService.GetPermissions(ctx)
		if err != nil {
			status := mapErrorStatus(err)
			res := response.BuildResponseFailed(fmt.Sprintf("%s permissions", dto.FAILED_GET_ALL), cleanErrorMessage(err))
			ctx.AbortWithStatusJSON(status, res)
			return
		}

		res := response.BuildResponseSuccess(fmt.Sprintf("%s permissions", dto.SUCCESS_GET_ALL), result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		ph.logger.Error("invalid get permissions query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s permissions", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ph.permissionService.GetPermissionsWithPagination(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s permissions", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s permissions", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}
