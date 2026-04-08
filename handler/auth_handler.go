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
