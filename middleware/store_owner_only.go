package middleware

import (
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func StoreOwnerOnly(
	storeRepo repository.IStoreRepository,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ambil store_id dari param
		storeIDStr := ctx.Param("store_id")
		storeID, err := uuid.Parse(storeIDStr)
		if err != nil {
			res := response.BuildResponseFailed(
				"failed authorization",
				dto.MESSAGE_FAILED_INVALID_UUID,
			)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		// ambil claims user login
		claimsRaw, exists := ctx.Get("claims")
		if !exists {
			res := response.BuildResponseFailed(
				"failed authorization",
				"unauthorized",
			)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		claims, ok := claimsRaw.(*jwt.CustomClaims)
		if !ok {
			res := response.BuildResponseFailed(
				"failed authorization",
				"invalid token",
			)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		// ambil store
		store, found, err := storeRepo.GetStoreByStoreID(
			ctx,
			nil,
			&storeID,
		)
		if err != nil {
			res := response.BuildResponseFailed(
				"failed authorization",
				dto.ErrInternal.Error(),
			)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		if !found {
			res := response.BuildResponseFailed(
				"failed authorization",
				dto.ErrNotFound.Error(),
			)
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}

		// cek owner
		if store.OwnerID == nil || store.OwnerID.String() != claims.UserID {
			res := response.BuildResponseFailed(
				"failed authorization",
				"only store owner can access this resource",
			)
			ctx.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}

		ctx.Next()
	}
}
