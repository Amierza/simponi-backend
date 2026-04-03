package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
)

func Vendor(route *gin.Engine, vendorHandler handler.IVendorHandler, jwtService jwt.IJWT, permissionService service.IPermissionService) {
	routes := route.Group("/api/v1/vendors").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/", middleware.RBAC(permissionService, "CreateVendor"), vendorHandler.CreateVendor)
		routes.GET("/", middleware.RBAC(permissionService, "GetVendors"), vendorHandler.GetVendors)
		routes.GET("/:id", middleware.RBAC(permissionService, "GetVendorByID"), vendorHandler.GetVendorByID)
		routes.PUT("/:id", middleware.RBAC(permissionService, "UpdateVendor"), vendorHandler.UpdateVendor)
		routes.DELETE("/:id", middleware.RBAC(permissionService, "DeleteVendorByID"), vendorHandler.DeleteVendorByID)
	}
}
