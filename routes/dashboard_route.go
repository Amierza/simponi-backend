package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Dashboard(route *gin.Engine, h handler.IDashboardHandler, jwtService jwt.IJWT) {
    routes := route.Group("/api/v1/stores/:store_id").Use(middleware.Authentication(jwtService))
    {
        routes.GET("/dashboard", h.GetDashboard)
    }
}
