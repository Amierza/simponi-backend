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
	IAuthHandler interface {
		SignUp(ctx *gin.Context)
		SignIn(ctx *gin.Context)
		RefreshToken(ctx *gin.Context)
	}

	authHandler struct {
		authService service.IAuthService
		logger      *zap.Logger
	}
)

func NewAuthHandler(authService service.IAuthService, logger *zap.Logger) *authHandler {
	return &authHandler{
		authService: authService,
		logger:      logger,
	}
}

// SignUp godoc
//
//	@Summary		Register new user
//	@Description	Register a new user with email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.SignUpRequest				true	"Sign up request"
//	@Success		201		{object}	dto.EmptySuccessResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse				"Bad Request - Invalid payload"
//	@Failure		500		{object}	dto.ErrorResponse				"Internal Server Error"
//	@Router			/auth/signup [post]
func (ah *authHandler) SignUp(ctx *gin.Context) {
	var payload dto.SignUpRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		ah.logger.Error("invalid signup request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s auth", dto.FAILED_SIGNUP), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	if err := ah.authService.SignUp(ctx, payload); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_SIGNUP, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_SIGNUP, nil)
	ctx.JSON(http.StatusCreated, res)
}

// SignIn godoc
//
//	@Summary		Sign in user
//	@Description	Authenticate user using email and password to obtain access token and refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.SignInRequest			true	"Sign in request"
//	@Success		200		{object}	dto.SignInResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse			"Bad Request - Invalid payload"
//	@Failure		401		{object}	dto.ErrorResponse			"Unauthorized - Invalid credentials"
//	@Failure		500		{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/auth/signin [post]
func (ah *authHandler) SignIn(ctx *gin.Context) {
	var payload dto.SignInRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		ah.logger.Error("invalid signin request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s auth", dto.FAILED_SIGNIN), cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ah.authService.SignIn(ctx, payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_SIGNIN, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_SIGNIN, result)
	ctx.JSON(http.StatusOK, res)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Generate new access token using valid refresh token (no Bearer token required)
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.RefreshTokenRequest			true	"Refresh token request"
//	@Success		200		{object}	dto.RefreshTokenResponseWrapper	"Success"
//	@Failure		400		{object}	dto.ErrorResponse				"Bad Request - Invalid payload"
//	@Failure		401		{object}	dto.ErrorResponse				"Unauthorized - Invalid refresh token"
//	@Failure		500		{object}	dto.ErrorResponse				"Internal Server Error"
//	@Router			/auth/refresh-token [post]
func (ah *authHandler) RefreshToken(ctx *gin.Context) {
	var payload dto.RefreshTokenRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		ah.logger.Error("invalid refresh token request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.MESSAGE_INVALID_REQUEST_PAYLOAD, err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := ah.authService.RefreshToken(ctx, payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_REFRESH_TOKEN, cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_REFRESH_TOKEN, result)
	ctx.JSON(http.StatusOK, res)
}
