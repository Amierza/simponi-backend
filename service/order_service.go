package service

import (
	"context"
	"fmt"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/jwt"
	"github.com/Amierza/simponi-backend/repository"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IOrderService interface {
		GetOrders(ctx context.Context, req response.PaginationRequest) (dto.OrderPaginationResponse, error)
		GetOrderByID(ctx context.Context, orderID *uuid.UUID) (dto.OrderResponse, error)
	}

	OrderService struct {
		orderRepo  repository.IOrderRepository
		logger     *zap.Logger
		jwtService jwt.IJWT
	}
)

func NewOrderService(orderRepo repository.IOrderRepository, logger *zap.Logger, jwtService jwt.IJWT) *OrderService {
	return &OrderService{
		orderRepo:  orderRepo,
		logger:     logger,
		jwtService: jwtService,
	}
}

func (os *OrderService) GetOrders(ctx context.Context, req response.PaginationRequest) (dto.OrderPaginationResponse, error) {
	datas, err := os.orderRepo.GetOrders(ctx, nil, &req)
	if err != nil {
		os.logger.Error("failed to get orders", zap.Error(err))
		return dto.OrderPaginationResponse{}, fmt.Errorf("failed to get orders: %w", dto.ErrInternal)
	}

	os.logger.Info("success to get orderes", zap.Int64("count", datas.Count))
	var orders []dto.OrderResponse
	for _, data := range datas.Orders {
		order := dto.OrderResponse{
			ID:               data.ID,
			ExternalOrderID:  data.ExternalOrderID,
			Ordernumber:      data.OrderNumber,
			StoreID:          data.StoreID,
			BuyerName:        data.BuyerName,
			BuyerPhone:       data.BuyerPhone,
			BuyerEmail:       data.BuyerEmail,
			ReceipentName:    data.RecipientName,
			ReceipentPhone:   data.RecipientPhone,
			ShippingAddress:  data.ShippingAddress,
			ShippingCity:     data.ShippingCity,
			ShippingProvince: data.ShippingProvince,
			Platform:         data.StorePlatform.Platform.Name,
			ShippingPostal:   data.ShippingPostal,
			ShippingMethod:   data.ShippingMethod,
			TrackingNumber:   data.TrackingNumber,
			SubtotalAmount:   data.SubtotalAmount,
			ShippingFee:      data.ShippingFee,
			MarketplaceFee:   data.MarketplaceFee,
			DiscountAmount:   data.DiscountAmount,
			TaxAmount:        data.TaxAmount,
			TotalAmount:      data.TotalAmount,
			NetAmount:        data.NetAmount,
			OrderStatus:      data.OrderStatus,
			PaymentStatus:    data.PaymentStatus,
			PaymentMethod:    data.PaymentMethod,
			OrderedAt:        data.OrderedAt,
			PaidAt:           data.PaidAt,
			ShippedAt:        data.ShippedAt,
			CompletedAt:      data.CompletedAt,
			CancelledAt:      data.CancelledAt,
			CreatedAt:        data.CreatedAt,
		}
		orders = append(orders, order)
	}

	return dto.OrderPaginationResponse{
		Data: orders,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: datas.MaxPage,
			Count:   datas.Count,
		},
	}, nil
}

func (os *OrderService) GetOrderByID(ctx context.Context, orderID *uuid.UUID) (dto.OrderResponse, error) {

	order, found, err := os.orderRepo.GetOrderByID(ctx, nil, orderID)

	if err != nil {
		os.logger.Error("failed to get order by ID", zap.String("order_id", orderID.String()), zap.Error(err))
		return dto.OrderResponse{}, fmt.Errorf("failed to get order by ID: %w", dto.ErrInternal)
	}

	if !found {
		os.logger.Warn("order not found", zap.String("order_id", orderID.String()))
		return dto.OrderResponse{}, fmt.Errorf("order not found: %w", dto.ErrNotFound)
	}

	os.logger.Info("success to get order by ID", zap.String("order_id", orderID.String()))

	return dto.OrderResponse{
		ID:               order.ID,
		ExternalOrderID:  order.ExternalOrderID,
		Ordernumber:      order.OrderNumber,
		StoreID:          order.StoreID,
		BuyerName:        order.BuyerName,
		BuyerPhone:       order.BuyerPhone,
		BuyerEmail:       order.BuyerEmail,
		ReceipentName:    order.RecipientName,
		ReceipentPhone:   order.RecipientPhone,
		ShippingAddress:  order.ShippingAddress,
		ShippingCity:     order.ShippingCity,
		ShippingProvince: order.ShippingProvince,
		ShippingPostal:   order.ShippingPostal,
		ShippingMethod:   order.ShippingMethod,
		TrackingNumber:   order.TrackingNumber,
		SubtotalAmount:   order.SubtotalAmount,
		ShippingFee:      order.ShippingFee,
		MarketplaceFee:   order.MarketplaceFee,
		DiscountAmount:   order.DiscountAmount,
		TaxAmount:        order.TaxAmount,
		TotalAmount:      order.TotalAmount,
		NetAmount:        order.NetAmount,
		OrderStatus:      order.OrderStatus,
		PaymentStatus:    order.PaymentStatus,
		PaymentMethod:    order.PaymentMethod,
		OrderedAt:        order.OrderedAt,
		PaidAt:           order.PaidAt,
		ShippedAt:        order.ShippedAt,
		CompletedAt:      order.CompletedAt,
		CancelledAt:      order.CancelledAt,
		CreatedAt:        order.CreatedAt,

		OrderDetails: []dto.OrderDetailResponse{},
	}, nil

}
