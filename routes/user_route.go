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
		routes.GET("/:user_id", middleware.RBAC(rolePermissionRepo, "GetUserByUserID"), userHandler.GetUserByUserID)
		routes.GET("/profile", middleware.RBAC(rolePermissionRepo, "GetUserProfile"), userHandler.GetUserProfile)

		routes.PUT("/:user_id", middleware.RBAC(rolePermissionRepo, "UpdateUserByUserID"), userHandler.UpdateUserByUserID)
		routes.PATCH("/:user_id/status", middleware.RBAC(rolePermissionRepo, "UpdateUserStatusByUserID"), userHandler.UpdateUserStatusByUserID)
		routes.PUT("/profile", middleware.RBAC(rolePermissionRepo, "UpdateUserProfile"), userHandler.UpdateUserProfile)

		routes.DELETE("/:user_id", middleware.RBAC(rolePermissionRepo, "DeleteUserByUserID"), userHandler.DeleteUserByUserID)
	}
}
