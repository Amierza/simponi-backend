// handler/platform_handler.go
package handler

import (
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	ISStorePlatformHandler interface {
		ConnectPlatform(*gin.Context)
		DisconnectPlatform(*gin.Context)
		SyncProducts(*gin.Context)
	}

	storePlatformHandler struct {
		storePlatformService service.IStorePlatformService
		logger               *zap.Logger
	}
)

func NewStorePlatformHandler(storePlatformService service.IStorePlatformService, logger *zap.Logger) *storePlatformHandler {
	return &storePlatformHandler{
		storePlatformService: storePlatformService,
		logger:               logger,
	}
}

// ConnectPlatform godoc
//
//	@Summary		Connect platform to store
//	@Description	Connect a platform to a store (Requires permission: ConnectPlatform)
//	@Tags			Stores
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_id		path	string	true	"Store ID"
//	@Param			platform_id		path	string	true	"Platform ID"
//	@Success		200				{object}	response.ResponseWrapper	"Success"
//	@Failure		400				{object}	response.ResponseWrapper	"Bad Request"
//	@Failure		401				{object}	response.ResponseWrapper	"Unauthorized"
//	@Failure		403				{object}	response.ResponseWrapper	"Forbidden"
//	@Failure		500				{object}	response.ResponseWrapper	"Internal Server Error"
//	@Router			/store/{store_id}/platforms/{platform_id}/connect [post]
func (sph *storePlatformHandler) ConnectPlatform(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		sph.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed("failed to connect platform", dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	platformIDStr := ctx.Param("platform_id")
	platformID, err := uuid.Parse(platformIDStr)
	if err != nil {
		sph.logger.Error("invalid platform ID", zap.String("platform_id", platformIDStr), zap.Error(err))
		res := response.BuildResponseFailed("failed to connect platform", dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sph.storePlatformService.ConnectPlatform(ctx, &storeID, &platformID)
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
//	@Summary		Disconnect platform from store
//	@Description	Disconnect a platform from a store (Requires permission: DisconnectPlatform)
//	@Tags			Stores
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_platform_id	path	string	true	"Store Platform ID"
//	@Success		200					{object}	response.ResponseWrapper	"Success"
//	@Failure		400					{object}	response.ResponseWrapper	"Bad Request"
//	@Failure		401					{object}	response.ResponseWrapper	"Unauthorized"
//	@Failure		403					{object}	response.ResponseWrapper	"Forbidden"
//	@Failure		500					{object}	response.ResponseWrapper	"Internal Server Error"
//	@Router			/store/{store_id}/platforms/{store_platform_id} [delete]
func (sph *storePlatformHandler) DisconnectPlatform(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		sph.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed("failed to disconnect platform", dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	platformIDStr := ctx.Param("platform_id")
	platformID, err := uuid.Parse(platformIDStr)
	if err != nil {
		sph.logger.Error("invalid platform ID", zap.String("platform_id", platformIDStr), zap.Error(err))
		res := response.BuildResponseFailed("failed to disconnect platform", dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = sph.storePlatformService.DisconnectPlatform(ctx, &storeID, &platformID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to disconnect platform", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("platform disconnected successfully", nil)
	ctx.JSON(http.StatusOK, res)
}

// SyncProducts godoc
//
//	@Summary		Sync products from platform to store
//	@Description	Sync products from a connected platform to the store (Requires permission: SyncPlatform)
//	@Tags			Stores
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			store_platform_id	path	string	true	"Store Platform ID"
//	@Success		200					{object}	response.ResponseWrapper	"Success"
//	@Failure		400					{object}	response.ResponseWrapper	"Bad Request"
//	@Failure		401					{object}	response.ResponseWrapper	"Unauthorized"
//	@Failure		403					{object}	response.ResponseWrapper	"Forbidden"
//	@Failure		500					{object}	response.ResponseWrapper	"Internal Server Error"
//	@Router			/store/{store_id}/platforms/{store_platform_id}/sync/products [post]
func (sph *storePlatformHandler) SyncProducts(ctx *gin.Context) {
	storeIDStr := ctx.Param("store_id")
	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		sph.logger.Error("invalid store ID", zap.String("store_id", storeIDStr), zap.Error(err))
		res := response.BuildResponseFailed("failed to sync products", dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	platformIDStr := ctx.Param("platform_id")
	platformID, err := uuid.Parse(platformIDStr)
	if err != nil {
		sph.logger.Error("invalid platform ID", zap.String("platform_id", platformIDStr), zap.Error(err))
		res := response.BuildResponseFailed("failed to sync products", dto.MESSAGE_FAILED_INVALID_UUID)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = sph.storePlatformService.SyncProducts(ctx, &storeID, &platformID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to sync products", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("products synced successfully", nil)
	ctx.JSON(http.StatusOK, res)
}
