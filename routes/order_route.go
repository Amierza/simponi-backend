package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Order(route *gin.Engine, orderHandler handler.IOrderHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/orders").Use(middleware.Authentication(jwtService))
	{
		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetOrders"), orderHandler.GetOrders)
		routes.GET("/:id", middleware.RBAC(rolePermissionRepo, "GetOrderByID"), orderHandler.GetOrderByID)
	}
}
