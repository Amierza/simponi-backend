// handler/platform_handler.go
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
	IPlatformHandler interface {
		GetMyStore(ctx *gin.Context)
		ConnectPlatform(ctx *gin.Context)
		DisconnectPlatform(ctx *gin.Context)
	}

	platformHandler struct {
		platformService service.IPlatformService
		logger          *zap.Logger
	}
)

func NewPlatformHandler(platformService service.IPlatformService, logger *zap.Logger) *platformHandler {
	return &platformHandler{
		platformService: platformService,
		logger:          logger,
	}
}

// GetMyStore godoc
//
//	@Summary		Get my store
//	@Description	Get store + connected platforms milik user yang sedang login
//	@Tags			Platforms
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{object}	dto.MyStoreResponse	"Store ditemukan"
//	@Success		204	"Belum punya store"
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/platforms/my-store [get]
func (ph *platformHandler) GetMyStore(ctx *gin.Context) {
	result, err := ph.platformService.GetMyStore(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to get store", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	// Belum punya store → 200 dengan data null, bukan 204
	// Supaya FE bisa parse JSON dengan konsisten
	res := response.BuildResponseSuccess("success get my store", result)
	ctx.JSON(http.StatusOK, res)
}

// ConnectPlatform godoc
//
//	@Summary		Connect platform (mock OAuth)
//	@Description	Connect Shopee/Tokopedia. Jika belum ada store, otomatis buat store baru.
//	@Tags			Platforms
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.ConnectPlatformRequest	true	"Connect request"
//	@Success		200		{object}	dto.MyStoreResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse	"Platform sudah terkoneksi"
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/platforms/connect [post]
func (ph *platformHandler) ConnectPlatform(ctx *gin.Context) {
	var payload dto.ConnectPlatformRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ph.logger.Error("invalid connect platform payload", zap.Error(err))
		res := response.BuildResponseFailed("failed to connect platform", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.platformService.ConnectPlatform(ctx, &payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to connect platform", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("platform connected successfully", result)
	ctx.JSON(http.StatusOK, res)
}

// DisconnectPlatform godoc
//
//	@Summary		Disconnect platform
//	@Description	Hapus koneksi platform. Jika ini platform terakhir, store juga dihapus.
//	@Tags			Platforms
//	@Security		BearerAuth
//	@Produce		json
//	@Param			store_platform_id	path	string	true	"StorePlatform ID (UUID)"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/platforms/{store_platform_id}/disconnect [delete]
func (ph *platformHandler) DisconnectPlatform(ctx *gin.Context) {
	idStr := ctx.Param("store_platform_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res := response.BuildResponseFailed(
			fmt.Sprintf("%s platform", dto.FAILED_DELETE),
			dto.MESSAGE_FAILED_INVALID_UUID,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if err := ph.platformService.DisconnectPlatform(ctx, &id); err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed(
			fmt.Sprintf("%s platform", dto.FAILED_DELETE),
			cleanErrorMessage(err),
		)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s platform", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}
