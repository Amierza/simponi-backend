package middleware

import (
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/gin-gonic/gin"
)

func RBAC(rolePermissionRepo repository.IRolePermissionRepository, permissionName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleIDRaw, exists := ctx.Get("role_id")
		if !exists {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_GET_ROLE_USER)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		roleID := roleIDRaw.(string)

		allowed, err := rolePermissionRepo.CheckRolePermissionByPermissionName(ctx, nil, roleID, permissionName)
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
