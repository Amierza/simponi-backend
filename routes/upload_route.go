package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Upload(route *gin.Engine, uploadHandler handler.IUploadHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/uploads").Use(middleware.Authentication(jwtService))
	{
		routes.POST("", uploadHandler.Upload)
	}
}
