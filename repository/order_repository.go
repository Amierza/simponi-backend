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
		CreateOrder(ctx context.Context, tx *gorm.DB, order *entity.Order) (*entity.Order, error)
		CreateOrderDetail(ctx context.Context, tx *gorm.DB, orderDetail *entity.OrderDetail) (*entity.OrderDetail, error)

		// READ
		GetOrders(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.OrderPaginationRepositoryResponse, error)
		GetOrderByID(ctx context.Context, tx *gorm.DB, orderID *uuid.UUID) (*entity.Order, bool, error)
		GetOrdersByVendorID(ctx context.Context, tx *gorm.DB, vendorID *uuid.UUID, req *response.PaginationRequest) (dto.OrderPaginationRepositoryResponse, error)
		GetOrdersByCustomerID(ctx context.Context, tx *gorm.DB, customerID *uuid.UUID, req *response.PaginationRequest) (dto.OrderPaginationRepositoryResponse, error)

		// UPDATE
		UpdateOrder(ctx context.Context, tx *gorm.DB, order *entity.Order) (*entity.Order, error)
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

func (or *OrderRepository) CreateOrder(ctx context.Context, tx *gorm.DB, order *entity.Order) (*entity.Order, error) {
	if tx == nil {
		tx = or.db
	}

	if err := tx.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil

}

func (or *OrderRepository) CreateOrderDetail(ctx context.Context, tx *gorm.DB, orderDetail *entity.OrderDetail) (*entity.OrderDetail, error) {
	if tx == nil {
		tx = or.db
	}

	if err := tx.WithContext(ctx).Create(orderDetail).Error; err != nil {
		return nil, err
	}
	return orderDetail, nil
}

func (or *OrderRepository) GetOrders(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.OrderPaginationRepositoryResponse, error) {
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

	query := tx.WithContext(ctx).Model(&entity.Order{}).Preload("OrderDetails").Preload("OrderDetails.Product").Preload("Vendor").Preload("Customer")

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("id::text ILIKE ? OR vendor_id::text ILIKE ? OR customer_id::text ILIKE ?", searchTerm, searchTerm, searchTerm)
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

func (or *OrderRepository) GetOrderByID(ctx context.Context, tx *gorm.DB, orderID *uuid.UUID) (*entity.Order, bool, error) {
	var order entity.Order
	if err := tx.WithContext(ctx).First(&order, "id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}
