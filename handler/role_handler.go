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
		GetRoleByID(ctx *gin.Context)
		UpdateRole(ctx *gin.Context)
		DeleteRoleByID(ctx *gin.Context)
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

func (rh *roleHandler) GetRoles(ctx *gin.Context) {
	result, err := rh.roleService.GetRoles(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s roles", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s roles", dto.SUCCESS_GET_ALL), result)
	ctx.JSON(http.StatusOK, res)
}

func (rh *roleHandler) GetRoleByID(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		rh.logger.Error("invalid role ID", zap.String("id", roleIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := rh.roleService.GetRoleByID(ctx, &roleID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (rh *roleHandler) UpdateRole(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
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

	result, err := rh.roleService.UpdateRole(ctx, &roleID, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (rh *roleHandler) DeleteRoleByID(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		rh.logger.Error("invalid role ID", zap.String("id", roleIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := rh.roleService.DeleteRoleByID(ctx, &roleID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s role", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s role", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
