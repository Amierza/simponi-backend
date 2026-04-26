package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Impersonate(route *gin.Engine, impersonateHandler handler.IImpersonateHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/impersonate").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/:user_id", middleware.RBAC(rolePermissionRepo, "StartImpersonate"), impersonateHandler.StartImpersonate)
		routes.POST("/stop", impersonateHandler.StopImpersonate)
	}
}
