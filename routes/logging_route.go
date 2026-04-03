package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
)

func Log(route *gin.Engine, logHandler handler.ILogHandler, jwtService jwt.IJWT, permissionService service.IPermissionService) {
	routes := route.Group("/api/v1/logs").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/logs", middleware.RBAC(permissionService, "CreateLog"), logHandler.CreateLog)
		routes.GET("/logs", middleware.RBAC(permissionService, "GetLogs"), logHandler.GetLogs)
		routes.GET("/logs/store/:storeID", middleware.RBAC(permissionService, "GetLogsByStoreID"), logHandler.GetLogsByStoreID)
		routes.GET("/logs/date-range", middleware.RBAC(permissionService, "GetLogsByDateRange"), logHandler.GetLogsByDateRange)
	}
}
