package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func ExternalProduct(route *gin.Engine, externalProductHandler handler.IExternalProductHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/external-products").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/product/:productId", middleware.RBAC(rolePermissionRepo, "GetExternalProductsByProductID"), externalProductHandler.GetExternalProductsByProductID)
		routes.GET("/store-platform/:storePlatformId", middleware.RBAC(rolePermissionRepo, "GetExternalProductsByStorePlatformID"), externalProductHandler.GetExternalProductsByStorePlatformID)
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateExternalProduct"), externalProductHandler.CreateExternalProduct)
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetExternalProducts"), externalProductHandler.GetExternalProducts)
		routes.GET("/:id", middleware.RBAC(rolePermissionRepo, "GetExternalProductByID"), externalProductHandler.GetExternalProductByID)
		routes.PUT("/:id", middleware.RBAC(rolePermissionRepo, "UpdateExternalProduct"), externalProductHandler.UpdateExternalProduct)
		routes.DELETE("/:id", middleware.RBAC(rolePermissionRepo, "DeleteExternalProductByID"), externalProductHandler.DeleteExternalProductByID)
	}
}
