package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"

	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Dashboard(route *gin.Engine, dashboardHandler handler.IDashboardHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/superadmin/dashboard").Use(middleware.Authentication(jwtService)).Use(middleware.OnlySuperAdmin())
	{
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetDashboardData"), dashboardHandler.GetDashboardData)
	}
}
