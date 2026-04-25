package middleware

import (
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/gin-gonic/gin"
)

func RBAC(rolePermissionRepo repository.IRolePermissionRepository, permissionName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. ambil claims dari context
		claimsRaw, exists := ctx.Get("claims")
		if !exists {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_GET_ROLE_USER)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		claims, ok := claimsRaw.(*jwt.CustomClaims)
		if !ok {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "invalid token claims")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		// 2. FAST CHECK (dari JWT)
		for _, p := range claims.Permissions {
			if p == permissionName {
				ctx.Next()
				return
			}
		}

		// 3. FALLBACK ke DB (optional, untuk safety)
		allowed, err := rolePermissionRepo.CheckRolePermissionByPermissionName(
			ctx,
			nil,
			claims.RoleID,
			permissionName,
		)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_CHECK_PERMISSION)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
			return
		}

		if !allowed {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_FORBIDDEN)
			ctx.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}

		ctx.Next()
	}
}
