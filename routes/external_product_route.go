package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func ExternalProduct(route *gin.Engine, externalProductHandler handler.IExternalProductHandler, jwtService jwt.IJWT, permissionService repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/external-products").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/product/:productId", middleware.RBAC(permissionService, "GetExternalProductsByProductID"), externalProductHandler.GetExternalProductsByProductID)
		routes.GET("/store-platform/:storePlatformId", middleware.RBAC(permissionService, "GetExternalProductsByStorePlatformID"), externalProductHandler.GetExternalProductsByStorePlatformID)
		routes.POST("/", middleware.RBAC(permissionService, "CreateExternalProduct"), externalProductHandler.CreateExternalProduct)
		routes.GET("/", middleware.RBAC(permissionService, "GetExternalProducts"), externalProductHandler.GetExternalProducts)
		routes.GET("/:id", middleware.RBAC(permissionService, "GetExternalProductByID"), externalProductHandler.GetExternalProductByID)
		routes.PUT("/:id", middleware.RBAC(permissionService, "UpdateExternalProduct"), externalProductHandler.UpdateExternalProduct)
		routes.DELETE("/:id", middleware.RBAC(permissionService, "DeleteExternalProductByID"), externalProductHandler.DeleteExternalProductByID)
	}
}