package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Product(route *gin.Engine, productHandler handler.IProductHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/products").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", productHandler.CreateProduct)
		routes.GET("/", productHandler.GetAllProducts)
		routes.GET("/:id", productHandler.GetProductByID)
		routes.GET("/sku", productHandler.GetProductBySKU)
		routes.GET("/category/:categoryId", productHandler.GetProductsByCategory)
		routes.PUT("/:id", productHandler.UpdateProduct)
		routes.PATCH("/:id/stock", productHandler.UpdateStock)
		routes.DELETE("/:id", productHandler.DeleteProduct)
	}
}