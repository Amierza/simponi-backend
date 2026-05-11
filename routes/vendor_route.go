package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/gin-gonic/gin"
)

func Vendor(route *gin.Engine, vendorHandler handler.IVendorHandler, jwtService jwt.IJWT, rolePermissionRepo repository.IRolePermissionRepository) {
	routes := route.Group("/api/v1/stores/:store_id/vendors").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(rolePermissionRepo, "CreateVendorByStoreID"), vendorHandler.CreateVendorByStoreID)

		routes.GET("/", middleware.RBAC(rolePermissionRepo, "GetVendorsByStoreID"), vendorHandler.GetVendorsByStoreID)
		routes.GET("/:vendor_id", middleware.RBAC(rolePermissionRepo, "GetVendorByStoreIDAndVendorID"), vendorHandler.GetVendorByStoreIDAndVendorID)
		// routes.GET("/:vendor_id/products", middleware.RBAC(rolePermissionRepo, "GetVendorProducts"), vendorHandler.GetVendorProducts)

		routes.PUT("/:vendor_id", middleware.RBAC(rolePermissionRepo, "UpdateVendorByStoreIDAndVendorID"), vendorHandler.UpdateVendorByStoreIDAndVendorID)

		routes.DELETE("/:vendor_id", middleware.RBAC(rolePermissionRepo, "DeleteVendorByStoreIDAndVendorID"), vendorHandler.DeleteVendorByStoreIDAndVendorID)
	}
}
