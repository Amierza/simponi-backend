package middleware

import (
	"net/http"

	"github.com/Amierza/simponi-backend/constants"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/response"
	"github.com/gin-gonic/gin"
)

func OnlySuperAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claimsRaw, exists := ctx.Get("claims")
		if !exists {
			res := response.BuildResponseFailed(
				"failed authorization",
				"claims not found",
			)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		claims, ok := claimsRaw.(*jwt.CustomClaims)
		if !ok {
			res := response.BuildResponseFailed(
				"failed authorization",
				"invalid claims type",
			)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if claims.RoleID != constants.SUPER_ADMIN_ROLE_ID {
			res := response.BuildResponseFailed(
				"failed authorization",
				"only super admin can access this resource",
			)
			ctx.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}

		ctx.Next()
	}
}
