package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Store(route *gin.Engine, storeHandler handler.IStoreHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/stores").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateStore"), storeHandler.CreateStore)

		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetStores"), storeHandler.GetStores)
		routes.GET("/:id", middleware.RBAC(rolePermissionRepo, "GetStoreByID"), storeHandler.GetStoreByID)

		routes.PUT("/:id", middleware.RBAC(rolePermissionRepo, "UpdateStore"), storeHandler.UpdateStore)

		routes.DELETE("/:id", middleware.RBAC(rolePermissionRepo, "DeleteStoreByID"), storeHandler.DeleteStoreByID)
	}
}
