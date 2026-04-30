package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func StoreUser(route *gin.Engine, storeUserHandler handler.IStoreUserHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/stores/:store_id/users").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateStoreUsers"), storeUserHandler.CreateStoreUsers)

		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetStoreUsers"), storeUserHandler.GetStoreUsers)
		routes.GET("/:user_id", middleware.RBAC(rolePermissionRepo, "GetStoreUserByStoreIDAndUserID"), storeUserHandler.GetStoreUserByStoreIDAndUserID)

		routes.DELETE("/:user_id", middleware.RBAC(rolePermissionRepo, "DeleteStoreUserByStoreIDAndUserID"), storeUserHandler.DeleteStoreUserByStoreIDAndUserID)
	}
}
