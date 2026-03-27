package routes

import (
	"github.com/Amierza/simponi-backend/handler"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Logging(route *gin.Engine, loggingHandler handler.ILoggingHandler, jwtService jwt.IJWT) {
	routes := route.Group("/api/v1/logs").Use(middleware.Authentication(jwtService))
	{
		routes.POST("/logs", loggingHandler.CreateLog)
		routes.GET("/logs", loggingHandler.GetLogs)
		routes.GET("/logs/store/:storeID", loggingHandler.GetLogsByStoreID)
		routes.GET("/logs/date-range", loggingHandler.GetLogsByDateRange)
	}
}
