package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/response"
	"github.com/Amierza/simponi-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IOrderHandler interface {
		GetOrders(ctx *gin.Context)
		GetOrderByID(ctx *gin.Context)
	}

	orderHandler struct {
		orderService service.IOrderService
		logger       *zap.Logger
	}
)

func NewOrderHandler(orderService service.IOrderService, logger *zap.Logger) *orderHandler {
	return &orderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

func (oh *orderHandler) GetOrders(ctx *gin.Context) {
	var payload response.PaginationRequest

	if err := ctx.ShouldBindQuery(&payload); err != nil {
		oh.logger.Error("invalid get orders request payload", zap.Error(err), zap.Any("payload", payload))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to get orders", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := oh.orderService.GetOrders(ctx, payload)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to get orders", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s orders", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (oh *orderHandler) GetOrderByID(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	log.Println(orderID)
	if err != nil {
		oh.logger.Error("invalid order ID format", zap.String("order_id", orderIDStr), zap.Error(err))
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to get order", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	result, err := oh.orderService.GetOrderByID(ctx, &orderID)
	if err != nil {
		status := mapErrorStatus(err)
		res := response.BuildResponseFailed("failed to get order", cleanErrorMessage(err))
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("success to get order", result)
	ctx.JSON(http.StatusOK, res)
}
