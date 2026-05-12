package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func OAuth(route *gin.Engine, oAuthHandler handler.IOAuthHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/oauth").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/:platform/callback", oAuthHandler.PlatformCallback)
	}
}
