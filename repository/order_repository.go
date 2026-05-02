package repository

import (
	"context"
	"math"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IOrderRepository interface {
		// CREATE
		// CreateOrder(ctx context.Context, tx *gorm.DB, order *entity.Order) (*entity.Order, error)
		// CreateOrderDetail(ctx context.Context, tx *gorm.DB, orderDetail *entity.OrderDetail) (*entity.OrderDetail, error)

		// READ
		GetOrders(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.OrderPaginationRepositoryResponse, error)
		GetOrderByID(ctx context.Context, tx *gorm.DB, orderID *uuid.UUID, storeID *uuid.UUID) (*entity.Order, bool, error)
		// GetOrdersByVendorID(ctx context.Context, tx *gorm.DB, vendorID *uuid.UUID, req *response.PaginationRequest) (dto.OrderPaginationRepositoryResponse, error)
		// GetOrdersByCustomerID(ctx context.Context, tx *gorm.DB, customerID *uuid.UUID, req *response.PaginationRequest) (dto.OrderPaginationRepositoryResponse, error)

		// UPDATE
		// UpdateOrder(ctx context.Context, tx *gorm.DB, order *entity.Order) (*entity.Order, error)
	}

	OrderRepository struct {
		db *gorm.DB
	}
)

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// func (or *OrderRepository) CreateOrder(ctx context.Context, tx *gorm.DB, order *entity.Order) (*entity.Order, error) {
// 	if tx == nil {
// 		tx = or.db
// 	}

// 	if err := tx.WithContext(ctx).Create(order).Error; err != nil {
// 		return nil, err
// 	}
// 	return order, nil

// }

// func (or *OrderRepository) CreateOrderDetail(ctx context.Context, tx *gorm.DB, orderDetail *entity.OrderDetail) (*entity.OrderDetail, error) {
// 	if tx == nil {
// 		tx = or.db
// 	}

// 	if err := tx.WithContext(ctx).Create(orderDetail).Error; err != nil {
// 		return nil, err
// 	}
// 	return orderDetail, nil
// }

func (or *OrderRepository) GetOrders(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, storeID *uuid.UUID) (dto.OrderPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = or.db
	}

	var orders []entity.Order
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Order{}).
		Where("store_id = ?", storeID).
		Preload("Store").
		Preload("StorePlatform").
		Preload("StorePlatform.Store").
		Preload("StorePlatform.Platform")

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("id::text ILIKE ? OR order_number::text ILIKE ? OR buyer_name::text ILIKE ? OR buyer_email::text ILIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.OrderPaginationRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(response.Paginate(req.Page, req.PerPage)).Find(&orders).Error; err != nil {
		return dto.OrderPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.OrderPaginationRepositoryResponse{
		Orders: orders,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, nil

}

func (or *OrderRepository) GetOrderByID(ctx context.Context, tx *gorm.DB, orderID *uuid.UUID, storeID *uuid.UUID) (*entity.Order, bool, error) {
	if tx == nil {
		tx = or.db
	}

	var order entity.Order
	if err := tx.WithContext(ctx).
		Preload("StorePlatform").
		Preload("StorePlatform.Store").
		Preload("StorePlatform.Platform").
		Preload("OrderDetails").
		Preload("OrderDetails.ExternalProduct").
		Preload("OrderDetails.ExternalProduct.Product").
		Preload("OrderDetails.ExternalProduct.Product.Images").
		Preload("OrderDetails.ExternalProduct.StorePlatform").
		Preload("OrderDetails.ExternalProduct.StorePlatform.Store").
		Preload("OrderDetails.ExternalProduct.StorePlatform.Platform").
		First(&order, "id = ? AND store_id = ?", orderID, storeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}
