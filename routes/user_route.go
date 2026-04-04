package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userHandler handler.IUserHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/users").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateUser"), userHandler.CreateUser)

		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetUsers"), userHandler.GetUsers)
		routes.GET("/:id", middleware.RBAC(rolePermissionRepo, "GetUserByID"), userHandler.GetUserByID)
		routes.GET("/profile", middleware.RBAC(rolePermissionRepo, "GetProfile"), userHandler.GetProfile)

		routes.PUT("/:id", middleware.RBAC(rolePermissionRepo, "UpdateUser"), userHandler.UpdateUser)

		routes.DELETE("/:id", middleware.RBAC(rolePermissionRepo, "DeleteUserByID"), userHandler.DeleteUserByID)
	}
}
