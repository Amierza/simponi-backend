package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userHandler handler.IUserHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/users").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/profile", userHandler.GetProfile)
	}
}
