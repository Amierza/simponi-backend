package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Vendor(route *gin.Engine, vendorHandler handler.IVendorHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/vendors").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateVendor"), vendorHandler.CreateVendor)

		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetVendors"), vendorHandler.GetVendors)
		routes.GET("/:id", middleware.RBAC(rolePermissionRepo, "GetVendorByID"), vendorHandler.GetVendorByID)

		routes.PUT("/:id", middleware.RBAC(rolePermissionRepo, "UpdateVendor"), vendorHandler.UpdateVendor)

		routes.DELETE("/:id", middleware.RBAC(rolePermissionRepo, "DeleteVendorByID"), vendorHandler.DeleteVendorByID)
	}
}
