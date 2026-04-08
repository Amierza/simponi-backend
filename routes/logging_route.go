package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Log(route *gin.Engine, logHandler handler.ILogHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/logs").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/logs", middleware.RBAC(rolePermissionRepo, "CreateLog"), logHandler.CreateLog)
		routes.GET("/logs", middleware.RBAC(rolePermissionRepo, "GetLogs"), logHandler.GetLogs)
		routes.GET("/logs/store/:storeID", middleware.RBAC(rolePermissionRepo, "GetLogsByStoreID"), logHandler.GetLogsByStoreID)
		routes.GET("/logs/date-range", middleware.RBAC(rolePermissionRepo, "GetLogsByDateRange"), logHandler.GetLogsByDateRange)
	}
}
