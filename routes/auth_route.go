package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/gin-gonic/gin"
)

func Auth(route *gin.Engine, authHandler handler.IAuthHandler) {
	routes := route.Group("/api/v1/auth")
	{
		routes.POST("/signin", authHandler.SignIn)
		routes.POST("/refresh-token", authHandler.RefreshToken)
	}
}
