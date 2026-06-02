package handler

import (
	"net/http"

	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type (
	IDashboardHandler interface {
		GetDashboardData(ctx *gin.Context)
	}

	dashboardHandler struct {
		dashboardService service.IDashboardService
		logger           *zap.Logger
	}
)

func NewDashboardHandler(dashboardService service.IDashboardService, logger *zap.Logger) *dashboardHandler {
	return &dashboardHandler{
		dashboardService: dashboardService,
		logger:           logger,
	}
}

// GetDashboardData godoc
//
//	@Summary		Get Dashboard Data
//	@Description	Get aggregated data for dashboard display (only for superadmin)
//	@Tags			Dashboard
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=dto.DashboardResponse}	"Success"
//	@Failure		400	{object}	response.Response								"Bad Request"
//	@Failure		401	{object}	response.Response								"Unauthorized"
//	@Failure		403	{object}	response.Response								"Forbidden"
//	@Failure		500	{object}	response.Response								"Internal Server Error"
//	@Security		ApiKeyAuth
//	@Router			/superadmin/dashboard [get]
func (dh *dashboardHandler) GetDashboardData(ctx *gin.Context) {
	data, err := dh.dashboardService.GetDashboardData(ctx)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to get dashboard data", err.Error())
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("success get dashboard data", data)
	ctx.JSON(http.StatusOK, res)

}
