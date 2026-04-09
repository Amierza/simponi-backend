package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func InventoryLog(route *gin.Engine, inventoryLogHandler handler.IInventoryLoggingHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/inventory-logs").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateInventoryLog"), inventoryLogHandler.CreateInventoryLog)
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetInventoryLogs"), inventoryLogHandler.GetInventoryLogs)
	}
}
