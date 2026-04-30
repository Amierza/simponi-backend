package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Product(route *gin.Engine, productHandler handler.IProductHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/stores/:store_id/products").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateProduct"), productHandler.CreateProduct)
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetProducts"), productHandler.GetProducts)
		routes.GET("/stats", middleware.RBAC(rolePermissionRepo, "GetProductStats"), productHandler.GetProductStats)
		routes.GET("/:product_id", middleware.RBAC(rolePermissionRepo, "GetProductByStoreIDAndProductID"), productHandler.GetProductByStoreIDAndProductID)
		routes.PUT("/:product_id", middleware.RBAC(rolePermissionRepo, "UpdateProductByStoreIDAndProductID"), productHandler.UpdateProductByStoreIDAndProductID)
		routes.PATCH("/:product_id/stock", middleware.RBAC(rolePermissionRepo, "UpdateStockByStoreIDAndProductID"), productHandler.UpdateStockByStoreIDAndProductID)
		routes.DELETE("/:product_id", middleware.RBAC(rolePermissionRepo, "DeleteProductByStoreIDAndProductID"), productHandler.DeleteProductByStoreIDAndProductID)
	}
}
