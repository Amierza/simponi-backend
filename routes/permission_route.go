package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Permission(route *gin.Engine, permissionHandler handler.IPermissionHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/permissions").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetPermissions"), permissionHandler.GetPermissions)
	}
}
