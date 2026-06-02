package repository

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/Amierza/simponi-backend/dto"
	"github.com/Amierza/simponi-backend/entity"
	"github.com/Amierza/simponi-backend/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	IStoreRepository interface {
		// CREATE
		CreateStore(ctx context.Context, tx *gorm.DB, store *entity.Store) error

		// READ
		GetStores(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.StorePaginationRepositoryResponse, error)
		GetStoresByUserID(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, userID *uuid.UUID) (dto.StorePaginationRepositoryResponse, error)
		GetStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (*entity.Store, bool, error)

		// UPDATE
		UpdateStoreByStoreID(ctx context.Context, tx *gorm.DB, store *entity.Store) error

		// DELETE
		DeleteStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) error

		// Dashboard
		GetSummary(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, from, to time.Time) (dto.DashboardSummaryMetricsResponse, error)
		GetTrend(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, months int) ([]dto.DashboardTrendPointResponse, error)
		GetRecentOrders(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, limit int) ([]dto.DashboardRecentOrderItemResponse, error)
		GetLowStock(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, threshold int, limit int) ([]dto.DashboardLowStockItemResponse, error)
		GetTopProducts(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, limit int) ([]dto.DashboardTopProductItemResponse, error)
		GetActivity(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, limit int) ([]dto.DashboardActivityItemResponse, error)
	}

	storeRepository struct {
		db *gorm.DB
	}
)

func NewStoreRepository(db *gorm.DB) *storeRepository {
	return &storeRepository{
		db: db,
	}
}

// CREATE
func (vr *storeRepository) CreateStore(ctx context.Context, tx *gorm.DB, store *entity.Store) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Create(store).Error
}

// READ
func (vr *storeRepository) GetStores(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest) (dto.StorePaginationRepositoryResponse, error) {
	if tx == nil {
		tx = vr.db
	}

	var stores []*entity.Store
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Preload("Owner").
		Preload("StorePlatforms.Platform").
		Preload("Orders").
		Preload("Logs")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.StorePaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&stores).Error; err != nil {
		return dto.StorePaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.StorePaginationRepositoryResponse{
		Stores: stores,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (vr *storeRepository) GetStoresByUserID(ctx context.Context, tx *gorm.DB, req *response.PaginationRequest, userID *uuid.UUID) (dto.StorePaginationRepositoryResponse, error) {
	if tx == nil {
		tx = vr.db
	}

	var stores []*entity.Store
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Joins("JOIN store_users su ON su.store_id = stores.id").
		Where("su.user_id = ?", userID).
		Preload("Owner").
		Preload("StorePlatforms.Platform").
		Preload("Orders").
		Preload("Logs")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.StorePaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(response.Paginate(req.Page, req.PerPage)).Find(&stores).Error; err != nil {
		return dto.StorePaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.StorePaginationRepositoryResponse{
		Stores: stores,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

func (vr *storeRepository) GetStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) (*entity.Store, bool, error) {
	if tx == nil {
		tx = vr.db
	}

	var store *entity.Store
	err := tx.WithContext(ctx).
		Model(&entity.Store{}).
		Preload("Owner").
		Preload("StorePlatforms.Platform").
		Preload("Orders").
		Preload("Logs").
		Where("id = ?", storeID).
		First(&store).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return store, true, nil
}

// UPDATE
func (vr *storeRepository) UpdateStoreByStoreID(ctx context.Context, tx *gorm.DB, store *entity.Store) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).
		Model(&entity.Store{}).
		Where("id = ?", store.ID).
		Updates(map[string]interface{}{
			"name":        store.Name,
			"description": store.Description,
			"image_url":   store.ImageURL,
			"is_active":   store.IsActive,
		}).Error
}

// DELETE
func (vr *storeRepository) DeleteStoreByStoreID(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID) error {
	if tx == nil {
		tx = vr.db
	}

	return tx.WithContext(ctx).Where("id = ?", &storeID).Delete(&entity.Store{}).Error
}

// Dashboard
func (vr *storeRepository) GetSummary(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, from, to time.Time) (dto.DashboardSummaryMetricsResponse, error) {
	db := vr.db
	if tx != nil {
		db = tx
	}
	var res dto.DashboardSummaryMetricsResponse

	// revenue and orders in range
	type revRes struct {
		Revenue int64
		Orders  int64
	}
	var r revRes
	if err := db.WithContext(ctx).
		Model(nil).
		Raw(`SELECT COALESCE(SUM(total_amount),0) AS revenue, COUNT(*) AS orders FROM orders WHERE store_id = ? AND ordered_at >= ? AND ordered_at <= ?`, storeID, from, to).
		Scan(&r).Error; err != nil {
		return res, err
	}
	res.RevenueMonthToDate = r.Revenue
	res.OrdersMonthToDate = r.Orders

	// products
	var cnt int64
	if err := db.WithContext(ctx).Model("products").Where("store_id = ?", *storeID).Count(&cnt).Error; err == nil {
		res.ActiveProducts = cnt
	}

	// low stock
	if err := db.WithContext(ctx).Model("products").Where("store_id = ? AND stock > 0 AND stock <= 10", *storeID).Count(&cnt).Error; err == nil {
		res.LowStockProducts = cnt
	}

	// out of stock
	if err := db.WithContext(ctx).Model("products").Where("store_id = ? AND stock = 0", *storeID).Count(&cnt).Error; err == nil {
		res.OutOfStockProducts = cnt
	}

	// orders by status
	var pending, ready, completed int64
	db.WithContext(ctx).Model("orders").Where("store_id = ? AND order_status = 'PENDING'", *storeID).Count(&pending)
	db.WithContext(ctx).Model("orders").Where("store_id = ? AND order_status = 'READY_TO_SHIP'", *storeID).Count(&ready)
	db.WithContext(ctx).Model("orders").Where("store_id = ? AND order_status = 'COMPLETED'", *storeID).Count(&completed)
	res.PendingOrders = pending
	res.ReadyToShipOrders = ready
	res.CompletedOrders = completed

	return res, nil
}

func (vr *storeRepository) GetTrend(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, months int) ([]dto.DashboardTrendPointResponse, error) {
	db := vr.db
	if tx != nil {
		db = tx
	}
	// simple approach: loop months and query sums per month
	var out []dto.DashboardTrendPointResponse
	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		start := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		type r struct {
			Revenue int64
			Orders  int64
		}
		var rr r
		if err := db.WithContext(ctx).Raw(`SELECT COALESCE(SUM(total_amount),0) AS revenue, COUNT(*) AS orders FROM orders WHERE store_id = ? AND ordered_at >= ? AND ordered_at <= ?`, storeID, start, end).Scan(&rr).Error; err != nil {
			return nil, err
		}
		out = append(out, dto.DashboardTrendPointResponse{Month: start.Format("2006-01"), Revenue: rr.Revenue, Orders: rr.Orders})
	}
	return out, nil
}

func (vr *storeRepository) GetRecentOrders(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, limit int) ([]dto.DashboardRecentOrderItemResponse, error) {
	db := vr.db
	if tx != nil {
		db = tx
	}
	var orders []struct {
		ID            uuid.UUID
		OrderNumber   string
		BuyerName     string
		OrderStatus   string
		PaymentStatus string
		TotalAmount   int64
		OrderedAt     *time.Time
	}
	if err := db.WithContext(ctx).Raw(`SELECT id, order_number, buyer_name, order_status, payment_status, total_amount, ordered_at FROM orders WHERE store_id = ? ORDER BY ordered_at DESC LIMIT ?`, storeID, limit).Scan(&orders).Error; err != nil {
		return nil, err
	}
	var out []dto.DashboardRecentOrderItemResponse
	for _, o := range orders {
		out = append(out, dto.DashboardRecentOrderItemResponse{ID: o.ID, OrderNumber: o.OrderNumber, BuyerName: o.BuyerName, OrderStatus: o.OrderStatus, PaymentStatus: o.PaymentStatus, TotalAmount: o.TotalAmount, OrderedAt: o.OrderedAt})
	}
	return out, nil
}

func (vr *storeRepository) GetLowStock(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, threshold int, limit int) ([]dto.DashboardLowStockItemResponse, error) {
	db := vr.db
	if tx != nil {
		db = tx
	}
	var products []struct {
		ID    uuid.UUID
		Name  string
		SKU   string
		Stock int
	}
	if err := db.WithContext(ctx).Raw(`SELECT id, name, sku, stock FROM products WHERE store_id = ? AND stock <= ? ORDER BY stock ASC LIMIT ?`, storeID, threshold, limit).Scan(&products).Error; err != nil {
		return nil, err
	}
	var out []dto.DashboardLowStockItemResponse
	for _, p := range products {
		out = append(out, dto.DashboardLowStockItemResponse{ProductID: p.ID, Name: p.Name, SKU: p.SKU, Stock: p.Stock, Threshold: threshold})
	}
	return out, nil
}

func (vr *storeRepository) GetTopProducts(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, limit int) ([]dto.DashboardTopProductItemResponse, error) {
	db := vr.db
	if tx != nil {
		db = tx
	}
	// join orders -> order_details -> external_products -> products
	var rows []struct {
		ProductID uuid.UUID
		Name      string
		SKU       string
		SoldQty   int64
		Revenue   int64
	}
	q := `SELECT p.id as product_id, p.name, p.sku, COALESCE(SUM(od.quantity),0) as sold_qty, COALESCE(SUM(od.quantity * ep.price),0) as revenue
        FROM order_details od
        JOIN orders o ON o.id = od.order_id
        JOIN external_products ep ON ep.id = od.external_product_id
        JOIN products p ON p.id = ep.product_id
        WHERE o.store_id = ?
        GROUP BY p.id, p.name, p.sku
        ORDER BY sold_qty DESC
        LIMIT ?`
	if err := db.WithContext(ctx).Raw(q, storeID, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	var out []dto.DashboardTopProductItemResponse
	for _, r := range rows {
		out = append(out, dto.DashboardTopProductItemResponse{ProductID: r.ProductID, Name: r.Name, SKU: r.SKU, SoldQty: r.SoldQty, Revenue: r.Revenue})
	}
	return out, nil
}

func (vr *storeRepository) GetActivity(ctx context.Context, tx *gorm.DB, storeID *uuid.UUID, limit int) ([]dto.DashboardActivityItemResponse, error) {
	db := vr.db
	if tx != nil {
		db = tx
	}
	var logs []struct {
		ID        uuid.UUID
		Action    string
		Message   string
		CreatedAt time.Time
	}
	if err := db.WithContext(ctx).Raw(`SELECT id, action, message, created_at FROM logs WHERE store_id = ? ORDER BY created_at DESC LIMIT ?`, storeID, limit).Scan(&logs).Error; err != nil {
		return nil, err
	}
	var out []dto.DashboardActivityItemResponse
	for _, l := range logs {
		out = append(out, dto.DashboardActivityItemResponse{ID: l.ID, Type: "log", Title: l.Action, Message: l.Message, CreatedAt: l.CreatedAt})
	}
	return out, nil
}
