package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Product(route *gin.Engine, productHandler handler.IProductHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/products").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/stats", middleware.RBAC(rolePermissionRepo, "GetProductStats"), productHandler.GetProductStats)
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateProduct"), productHandler.CreateProduct)
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetAllProducts"), productHandler.GetProducts)
		routes.GET("/:id", middleware.RBAC(rolePermissionRepo, "GetProductByID"), productHandler.GetProductByID)
		routes.GET("/sku", middleware.RBAC(rolePermissionRepo, "GetProductBySKU"), productHandler.GetProductBySKU)
		routes.GET("/category/:categoryId", middleware.RBAC(rolePermissionRepo, "GetProductsByCategory"), productHandler.GetProductsByCategoryID)
		routes.PUT("/:id", middleware.RBAC(rolePermissionRepo, "UpdateProduct"), productHandler.UpdateProduct)
		routes.PATCH("/:id/stock", middleware.RBAC(rolePermissionRepo, "UpdateStock"), productHandler.UpdateStock)
		routes.DELETE("/:id", middleware.RBAC(rolePermissionRepo, "DeleteProduct"), productHandler.DeleteProductByID)
	}
}
