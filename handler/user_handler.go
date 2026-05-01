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
		GetUserByUserID(ctx *gin.Context)
		GetUserProfile(ctx *gin.Context)
		UpdateUserByUserID(ctx *gin.Context)
		UpdateUserStatusByUserID(ctx *gin.Context)
		UpdateUserProfile(ctx *gin.Context)
		DeleteUserByUserID(ctx *gin.Context)
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

// CreateUser godoc
//
//	@Summary		Create new user
//	@Description	Create a new user (Requires permission: CreateUser)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateUserRequest	true	"Create user request"
//	@Success		201		{object}	dto.UserResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Bad Request"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden - Insufficient permission"
//	@Failure		409		{object}	dto.ErrorResponse		"Conflict - User already exists"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/users [post]
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

// GetUsers godoc
//
//	@Summary		Get list of users
//	@Description	Get paginated users (Requires permission: GetUsers)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int							false	"Page number"
//	@Param			limit	query		int							false	"Items per page"
//	@Success		200		{object}	dto.UsersResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse			"Bad Request"
//	@Failure		401		{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse			"Forbidden"
//	@Failure		500		{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/users [get]
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

// GetUserByUserID godoc
//
//	@Summary		Get user by ID
//	@Description	Get user detail by user ID (Requires permission: GetUserByUserID)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string					true	"User ID (UUID)"
//	@Success		200		{object}	dto.UserResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid UUID"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse		"User not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/users/{user_id} [get]
func (uh *userHandler) GetUserByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := uh.userService.GetUserByUserID(ctx, &userID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_DETAIL), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

// GetUserProfile godoc
//
//	@Summary		Get current user profile
//	@Description	Get authenticated user profile (Requires permission: GetUserProfile)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UserResponseWrapper	"Success"
//	@Failure		401	{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		500	{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/users/profile [get]
func (uh *userHandler) GetUserProfile(ctx *gin.Context) {
	result, err := uh.userService.GetUserProfile(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_GET_PROFILE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_GET_PROFILE), result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateUserByUserID godoc
//
//	@Summary		Update user
//	@Description	Update user data by ID (Requires permission: UpdateUserByUserID)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string					true	"User ID (UUID)"
//	@Param			payload	body		dto.UpdateUserRequest	true	"Update user request"
//	@Success		200		{object}	dto.UserResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid input"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse		"User not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/users/{user_id} [put]
func (uh *userHandler) UpdateUserByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateUserRequest
	payload.ID = userID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		uh.logger.Error("invalid update user request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.UpdateUserByUserID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateUserStatusByUserID godoc
//
//	@Summary		Update user status
//	@Description	Update user status (active/inactive) (Requires permission: UpdateUserStatusByUserID)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string					true	"User ID (UUID)"
//	@Param			payload	body		dto.UpdateUserStatus	true	"Update status request"
//	@Success		200		{object}	dto.UserResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid input"
//	@Failure		401		{object}	dto.ErrorResponse		"Unauthorized"
//	@Failure		403		{object}	dto.ErrorResponse		"Forbidden"
//	@Failure		404		{object}	dto.ErrorResponse		"User not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/users/{user_id}/status [patch]
func (uh *userHandler) UpdateUserStatusByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_UPDATE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var payload dto.UpdateUserStatus
	payload.ID = userID
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		uh.logger.Error("invalid update user status request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user status", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.UpdateUserStatusByUserID(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user status", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user status", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateUserProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update current logged-in user profile
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		dto.UpdateUserRequest	true	"Update user profile request"
//	@Success		200		{object}	response.Response{data=dto.UserResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/users/profile [put]
func (uh *userHandler) UpdateUserProfile(ctx *gin.Context) {
	var payload dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		uh.logger.Error("invalid update user request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user profile", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := uh.userService.UpdateUserProfile(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user profile", dto.FAILED_UPDATE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user profile", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

// DeleteUserByUserID godoc
//
//	@Summary		Delete user
//	@Description	Delete user by ID (Requires permission: DeleteUserByUserID)
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UserEmptyResponseWrapper	"Success"
//	@Failure		400	{object}	dto.ErrorResponse				"Invalid UUID"
//	@Failure		401	{object}	dto.ErrorResponse				"Unauthorized"
//	@Failure		403	{object}	dto.ErrorResponse				"Forbidden"
//	@Failure		404	{object}	dto.ErrorResponse				"User not found"
//	@Failure		500	{object}	dto.ErrorResponse				"Internal Server Error"
//	@Router			/users/{user_id} [delete]
func (uh *userHandler) DeleteUserByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		uh.logger.Error("invalid user ID", zap.String("id", userIDStr), zap.Error(err))
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_DELETE), dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := uh.userService.DeleteUserByUserID(ctx, &userID); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s user", dto.FAILED_DELETE), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s user", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
