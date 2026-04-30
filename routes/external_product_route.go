package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func ExternalProduct(route *gin.Engine, externalProductHandler handler.IExternalProductHandler, jwtService jwt.IJWT, permissionService repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/stores/:store_id/external-products").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(permissionService, "CreateExternalProduct"), externalProductHandler.CreateExternalProduct)
		routes.GET("/", middleware.RBAC(permissionService, "GetExternalProducts"), externalProductHandler.GetExternalProducts)
		routes.GET("/:external_product_id", middleware.RBAC(permissionService, "GetExternalProductByStoreIDAndExprodID"), externalProductHandler.GetExternalProductByStoreIDAndExprodID)
		routes.GET("/store-platform/:store_platform_id", middleware.RBAC(permissionService, "GetExternalProductsByStoreIDAndStorePlatformID"), externalProductHandler.GetExternalProductsByStoreIDAndStorePlatformID)
		routes.PUT("/:external_product_id", middleware.RBAC(permissionService, "UpdateExternalProductByStoreIDAndExprodID"), externalProductHandler.UpdateExternalProductByStoreIDAndExprodID)
		routes.DELETE("/:external_product_id", middleware.RBAC(permissionService, "DeleteExternalProductByStoreIDAndExprodID"), externalProductHandler.DeleteExternalProductByStoreIDAndExprodID)
	}
}
