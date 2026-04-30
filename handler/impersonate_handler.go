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

// StartImpersonate godoc
//
//	@Summary		Start impersonation
//	@Description	Admin impersonates another user and receives a new access token
//	@Tags			Impersonate
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			user_id	path		string	true	"User ID"
//	@Success		200		{object}	dto.ImpersonateSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Router			/impersonate/{user_id} [post]
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

// StopImpersonate godoc
//
//	@Summary		Stop impersonation
//	@Description	Stop impersonation and return to original admin user
//	@Tags			Impersonate
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	dto.ImpersonateSuccessResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/impersonate/stop [post]
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
