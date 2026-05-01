package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/gin-gonic/gin"
)

func ProductCategories(route *gin.Engine, productCategoriesHandler handler.IProductCategoriesHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/products/categories")
	{
		routes.GET("/", productCategoriesHandler.GetProductCategories)
	}
}
