package handler

import (
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IAuthHandler interface {
		SignIn(ctx *gin.Context)
		RefreshToken(ctx *gin.Context)
	}

	authHandler struct {
		authService service.IAuthService
	}
)

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (ah *authHandler) SignIn(ctx *gin.Context){
	var payload dto.SignInRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_INVALID_REQUEST_PAYLOAD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.authService.SignIn(ctx, payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_SIGNIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_SIGNIN, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *authHandler) RefreshToken(ctx *gin.Context){
	var payload dto.RefreshTokenRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_INVALID_REQUEST_PAYLOAD, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.authService.RefreshToken(ctx, payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_REFRESH_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_REFRESH_TOKEN, result)
	ctx.JSON(http.StatusOK, res)
}

