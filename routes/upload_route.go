package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Upload(route *gin.Engine, uploadHandler handler.IUploadHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/uploads").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "Upload"), uploadHandler.Upload)
	}
}
