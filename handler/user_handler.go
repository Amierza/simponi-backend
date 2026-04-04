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
	IUserHandler interface {
		CreateUser(ctx *gin.Context)
		GetUsers(ctx *gin.Context)
		GetUserByID(ctx *gin.Context)
		GetProfile(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
		DeleteUserByID(ctx *gin.Context)
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

func (uh *userHandler) CreateUser(ctx *gin.Context) {
	var payload dto.CreateUserRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		uh.logger.Error("invalid create user request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.CreateUser(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_CREATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusCreated, res)
}

func (uh *userHandler) GetUsers(ctx *gin.Context) {
	var payload response.PaginationRequest
	if err := ctx.ShouldBindQuery(&payload); err != nil {
		uh.logger.Error("invalid get users query payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s users", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.GetUsers(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s users", dto.FAILED_GET_ALL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s users", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) GetUserByID(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.GetUserByID(ctx, &userID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) GetProfile(ctx *gin.Context) {
	userIDAny, exists := ctx.Get("user_id")
	if !exists {
		res := response.BuildResponseFailed(dto.UNAUTHORIZED, dto.UNAUTHORIZED)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	userIDStr := userIDAny.(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.GetProfile(ctx, &userID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_PROFILE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_GET_PROFILE), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) UpdateUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		uh.logger.Error("invalid update user request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.UpdateUser(ctx, &userID, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (uh *userHandler) DeleteUserByID(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := uh.userService.DeleteUserByID(ctx, &userID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
