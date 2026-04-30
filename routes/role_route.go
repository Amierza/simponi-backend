package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Role(route *gin.Engine, roleHandler handler.IRoleHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/roles").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateRole"), roleHandler.CreateRole)

		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetRoles"), roleHandler.GetRoles)
		routes.GET("/:role_id", middleware.RBAC(rolePermissionRepo, "GetRoleByRoleID"), roleHandler.GetRoleByRoleID)

		routes.PUT("/:role_id", middleware.RBAC(rolePermissionRepo, "UpdateRoleByRoleID"), roleHandler.UpdateRoleByRoleID)

		routes.DELETE("/:role_id", middleware.RBAC(rolePermissionRepo, "DeleteRoleByRoleID"), roleHandler.DeleteRoleByRoleID)
	}
}
