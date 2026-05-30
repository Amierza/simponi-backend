package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
    IDashboardHandler interface {
        GetDashboard(ctx *gin.Context)
    }

    dashboardHandler struct {
        dashboardService service.IDashboardService
        storeUserService service.IStoreUserService
        logger *zap.Logger
    }
)

func NewDashboardHandler(dashboardService service.IDashboardService, storeUserService service.IStoreUserService, logger *zap.Logger) *dashboardHandler {
    return &dashboardHandler{dashboardService: dashboardService, storeUserService: storeUserService, logger: logger}
}

// GetDashboard godoc
// @Summary Get dashboard
// @Description Get consolidated dashboard for a store
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param store_id path string true "Store ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/stores/{store_id}/dashboard [get]
func (dh *dashboardHandler) GetDashboard(ctx *gin.Context) {
    storeIDStr := ctx.Param("store_id")
    storeID, err := uuid.Parse(storeIDStr)
    if err != nil {
        dh.logger.Error("invalid store ID", zap.String("id", storeIDStr), zap.Error(err))
        res := response.BuildResponseFailed(fmt.Sprintf("%s dashboard", dto.FAILED_GET_DETAIL), dto.MESSAGE_FAILED_INVALID_UUID)
        ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
        return
    }

    // check membership
    userIDVal, exists := ctx.Get("user_id")
    if !exists {
        res := response.BuildResponseFailed("failed to get dashboard", "missing user context")
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
        return
    }
    userID, ok := userIDVal.(uuid.UUID)
    if !ok {
        // sometimes stored as string
        userIDStr, okStr := userIDVal.(string)
        if !okStr {
            res := response.BuildResponseFailed("failed to get dashboard", "invalid user context")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
            return
        }
        userID, _ = uuid.Parse(userIDStr)
    }

    _, err = dh.storeUserService.GetStoreUserByStoreIDAndUserID(ctx, &storeID, &userID)
    if err != nil {
        if errors.Is(err, dto.ErrNotFound) {
            res := response.BuildResponseFailed("unauthorized to access store dashboard", "user not member of store")
            ctx.AbortWithStatusJSON(http.StatusForbidden, res)
            return
        }
        dh.logger.Error("failed to check store membership", zap.Error(err))
        res := response.BuildResponseFailed("failed to get dashboard", "internal error")
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
        return
    }

    result, err := dh.dashboardService.GetDashboard(ctx, &storeID)
    if err != nil {
        dh.logger.Error("failed to get dashboard data", zap.Error(err))
        res := response.BuildResponseFailed("failed to get dashboard", "internal error")
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
        return
    }

    res := response.BuildResponseSuccess("success to get dashboard", result)
    ctx.JSON(http.StatusOK, res)
}
