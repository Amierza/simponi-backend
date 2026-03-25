package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/gin-gonic/gin"
)

func Auth(route *gin.Engine, authHandler handler.IAuthHandler, jwtService jwt.IJWT){
	routes := route.Group("/api/v1/auth")
	{
		routes.POST("/signin", authHandler.SignIn)
		routes.POST("/refresh-token", authHandler.RefreshToken)
	}
}