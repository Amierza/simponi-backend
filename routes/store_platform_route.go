// routes/platform_route.go
package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func StorePlatform(route *gin.Engine, storePlatformHandler handler.ISStorePlatformHandler, rolePermissionRepo repository.IRolePermissionRepository, jwtService jwt.IJWT) {
	base := route.Group("/api/v1/store/:store_id/platforms")
	base.Use(middleware.Authentication(jwtService))
	{
		// connect platform (generate auth url)
		base.POST("/:platform_id/connect", middleware.RBAC(rolePermissionRepo, "ConnectPlatform"), storePlatformHandler.ConnectPlatform)

		// disconnect platform
		base.DELETE("/:platform_id", middleware.RBAC(rolePermissionRepo, "DisconnectPlatform"), storePlatformHandler.DisconnectPlatform)

		// sync routes
		syncRoutes := base.Group("/:platform_id/sync", middleware.RBAC(rolePermissionRepo, "SyncPlatform"))
		{
			syncRoutes.POST("/products", storePlatformHandler.SyncProducts)
			// syncRoutes.POST("/orders", storePlatformHandler.SyncOrders)
		}
	}
}
