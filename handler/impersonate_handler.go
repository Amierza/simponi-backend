package handler

import (
	"net/http"

	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	IImpersonateHandler interface {
		StartImpersonate(ctx *gin.Context)
		StopImpersonate(ctx *gin.Context)
	}

	impersonateHandler struct {
		impersonateService service.IImpersonateService
		logger             *zap.Logger
	}
)

func NewImpersonateHandler(impersonateService service.IImpersonateService, logger *zap.Logger) *impersonateHandler {
	return &impersonateHandler{
		impersonateService: impersonateService,
		logger:             logger,
	}
}

func (ih *impersonateHandler) StartImpersonate(ctx *gin.Context) {
	targetUserID := ctx.Param("user_id")

	// ambil admin dari claims
	claimsRaw, exists := ctx.Get("claims")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildResponseFailed("failed", "unauthorized"))
		return
	}

	claims, ok := claimsRaw.(*jwt.CustomClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildResponseFailed("failed", "invalid token"))
		return
	}

	result, err := ih.impersonateService.StartImpersonate(ctx, claims.UserID, targetUserID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed impersonate", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("success impersonate", result)
	ctx.JSON(http.StatusOK, res)
}

func (ih *impersonateHandler) StopImpersonate(ctx *gin.Context) {
	claimsRaw, exists := ctx.Get("claims")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildResponseFailed("failed", "unauthorized"))
		return
	}

	claims, ok := claimsRaw.(*jwt.CustomClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.BuildResponseFailed("failed", "invalid token"))
		return
	}

	result, err := ih.impersonateService.StopImpersonate(ctx, claims)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed stop impersonate", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("success stop impersonate", result)
	ctx.JSON(http.StatusOK, res)
}
