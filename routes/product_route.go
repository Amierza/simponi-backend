package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
)

func Product(route *gin.Engine, productHandler handler.IProductHandler, jwtService jwt.IJWT, permissionService service.IPermissionService) {
	routes := route.Group("/api/v1/products").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/stats", middleware.RBAC(permissionService, "GetProductStats"), productHandler.GetProductStats)
		routes.GET("/sku", middleware.RBAC(permissionService, "GetProductBySKU"), productHandler.GetProductBySKU)
		routes.GET("/category/:categoryId", middleware.RBAC(permissionService, "GetProductsByCategoryID"), productHandler.GetProductsByCategoryID)
		routes.POST("/", middleware.RBAC(permissionService, "CreateProduct"), productHandler.CreateProduct)
		routes.GET("/", middleware.RBAC(permissionService, "GetProducts"), productHandler.GetProducts)
		routes.GET("/:id", middleware.RBAC(permissionService, "GetProductByID"), productHandler.GetProductByID)
		routes.PUT("/:id", middleware.RBAC(permissionService, "UpdateProduct"), productHandler.UpdateProduct)
		routes.PATCH("/:id/stock", middleware.RBAC(permissionService, "UpdateStock"), productHandler.UpdateStock)
		routes.DELETE("/:id", middleware.RBAC(permissionService, "DeleteProductByID"), productHandler.DeleteProductByID)
	}
}