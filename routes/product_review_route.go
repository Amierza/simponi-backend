package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func ProductReview(route *gin.Engine, productReviewHandler handler.IProductReviewHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/stores/:store_id/products/:product_id/reviews").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateProductReview"), productReviewHandler.CreateProductReview)
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetProductReviews"), productReviewHandler.GetProductReviews)
	}
}
