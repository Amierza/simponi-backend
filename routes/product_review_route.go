package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func ProductReview(route *gin.Engine, productReviewHandler handler.IProductReviewHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	// create review is scoped to a specific product
	productScoped := route.Group("/api/v1/stores/:store_id/products/:product_id/reviews").Use(middleware.Authentication(jwtService))
	{
		productScoped.POST("/", middleware.RBAC(rolePermissionRepo, "CreateProductReview"), productReviewHandler.CreateProductReview)
	}

	// listing reviews is scoped to the whole store (all products in the store)
	storeScoped := route.Group("/api/v1/stores/:store_id/reviews").Use(middleware.Authentication(jwtService))
	{
		storeScoped.GET("/", middleware.RBAC(rolePermissionRepo, "GetProductReviews"), productReviewHandler.GetProductReviewsByStoreID)
	}
}
