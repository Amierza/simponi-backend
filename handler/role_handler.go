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
	IRoleHandler interface {
		CreateRole(ctx *gin.Context)
		GetRoles(ctx *gin.Context)
		GetRoleByRoleID(ctx *gin.Context)
		UpdateRoleByRoleID(ctx *gin.Context)
		DeleteRoleByRoleID(ctx *gin.Context)
	}

	roleHandler struct {
		roleService service.IRoleService
		logger      *zap.Logger
	}
)

func NewRoleHandler(roleService service.IRoleService, logger *zap.Logger) *roleHandler {
	return &roleHandler{
		roleService: roleService,
		logger:      logger,
	}
}

// CreateRole godoc
//
//	@Summary		Create new role
//	@Description	Create a new role with permissions (Requires permission: CreateRole)
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateRoleRequest	true	"Create role request"
//	@Success		201		{object}	dto.RoleResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Bad Request"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		409		{object}	dto.ErrorResponse		"Conflict - Role already exists"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/roles [post]
func (rh *roleHandler) CreateRole(ctx *gin.Context) {
	var payload dto.CreateRoleRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		rh.logger.Error("invalid create role request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := rh.roleService.CreateRole(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

// GetRoles godoc
//
//	@Summary		Get list of roles
//	@Description	Get paginated roles (Requires permission: GetRoles)
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int							false	"Page number"
//	@Param			limit	query		int							false	"Items per page"
//	@Success		200		{object}	dto.RolesResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse			"Bad Request"
//	@Failure		401		{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse			"Forbidden"
//	@Failure		500		{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/roles [get]
func (rh *roleHandler) GetRoles(ctx *gin.Context) {
	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		rh.logger.Error("invalid get roles query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s roles", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := rh.roleService.GetRoles(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s roles", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s roles", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

// GetRoleByRoleID godoc
//
//	@Summary		Get role by ID
//	@Description	Get role detail with permissions (Requires permission: GetRoleByRoleID)
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			role_id	path		string					true	"Role ID (UUID)"
//	@Success		200		{object}	dto.RoleResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid UUID"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse		"Role not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/roles/{role_id} [get]
func (rh *roleHandler) GetRoleByRoleID(ctx *gin.Context) {
	roleIDStr := ctx.Param("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		rh.logger.Error("invalid role ID", zap.String("id", roleIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := rh.roleService.GetRoleByRoleID(ctx, &roleID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateRoleByRoleID godoc
//
//	@Summary		Update role
//	@Description	Update role and its permissions (Requires permission: UpdateRoleByRoleID)
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			role_id	path		string					true	"Role ID (UUID)"
//	@Param			payload	body		dto.UpdateRoleRequest	true	"Update role request"
//	@Success		200		{object}	dto.RoleResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid input"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse		"Role not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/roles/{role_id} [put]
func (rh *roleHandler) UpdateRoleByRoleID(ctx *gin.Context) {
	roleIDStr := ctx.Param("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		rh.logger.Error("invalid role ID", zap.String("id", roleIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		rh.logger.Error("invalid update role request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := rh.roleService.UpdateRoleByRoleID(ctx, &roleID, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

// DeleteRoleByRoleID godoc
//
//	@Summary		Delete role
//	@Description	Delete role by ID (Requires permission: DeleteRoleByRoleID)
//	@Tags			Roles
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.RoleEmptyResponseWrapper	"Success"
//	@Failure		400	{object}	dto.ErrorResponse				"Invalid UUID"
//	@Failure		401	{object}	dto.ErrorResponse				"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse				"Forbidden"
//	@Failure		404	{object}	dto.ErrorResponse				"Role not found"
//	@Failure		500	{object}	dto.ErrorResponse				"Internal Server Error"
//	@Router			/roles/{role_id} [delete]
func (rh *roleHandler) DeleteRoleByRoleID(ctx *gin.Context) {
	roleIDStr := ctx.Param("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		rh.logger.Error("invalid role ID", zap.String("id", roleIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := rh.roleService.DeleteRoleByRoleID(ctx, &roleID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
