package handler

import (
	"net/http"

	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type (
	IOAuthHandler interface {
		PlatformCallback(*gin.Context)
	}

	oAuthHandler struct {
		oAuthService service.IOAuthService
		logger       *zap.Logger
	}
)

func NewOAuthHandler(oAuthService service.IOAuthService, logger *zap.Logger) *oAuthHandler {
	return &oAuthHandler{
		oAuthService: oAuthService,
		logger:       logger,
	}
}

func (oah *oAuthHandler) PlatformCallback(ctx *gin.Context) {
	platform := ctx.Param("platform")
	oah.logger.Info("received oauth callback", zap.String("platform", platform))

	var err error

	switch platform {
	case "shopee":
		code := ctx.Query("code")
		shopID := ctx.Query("shop_id")
		state := ctx.Query("state")
		err = oah.oAuthService.HandleShopeeCallback(ctx, code, shopID, state)
	default:
		oah.logger.Warn("unsupported platform for oauth callback", zap.String("platform", platform))
		ctx.JSON(400, gin.H{"error": "unsupported platform"})
	}

	if err != nil {
		oah.logger.Error("failed to handle oauth callback", zap.String("platform", platform), zap.Error(err))
		res := response.BuildResponseFailed("failed to handle oauth callback", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponseSuccess("oauth callback handled successfully", nil)
	ctx.JSON(http.StatusOK, res)
}
