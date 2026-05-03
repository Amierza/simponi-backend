// routes/platform_route.go
package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/gin-gonic/gin"
)

// Platform routes tidak menggunakan RBAC per-endpoint
// karena ketiga endpoint ini hanya beroperasi pada data milik user yang login.
// Ownership divalidasi di dalam service layer.
func Platform(
	route *gin.Engine,
	platformHandler handler.IPlatformHandler,
	jwtService jwt.IJWT,
) {
	routes := route.Group("/api/v1/platforms").Use(middleware.Authentication(jwtService))
	{
		// GET /platforms/my-store → cek status koneksi user
		routes.GET("/my-store", platformHandler.GetMyStore)

		// POST /platforms/connect → connect platform (mock OAuth)
		routes.POST("/connect", platformHandler.ConnectPlatform)

		// DELETE /platforms/:store_platform_id/disconnect → disconnect platform
		routes.DELETE("/:store_platform_id/disconnect", platformHandler.DisconnectPlatform)
	}
}
