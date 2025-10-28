package repositories

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/common"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	GetNumbering(ctx context.Context, module string) (*models.Numbering, error)
	CreateOrder(ctx context.Context, order *models.Order) error
	UpdateOrder(ctx context.Context, order *models.Order) error
	UpdateNumbering(ctx context.Context, numbering *models.Numbering) error
	CreateOrderItem(ctx context.Context, orderItem *models.OrderItem) error
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	GetListOrderAdmin(ctx context.Context, pagination *common.PaginationRequest) ([]*models.Order, *common.PaginationResponse, error)
	GetListOrder(ctx context.Context, userID uint, pagination *common.PaginationRequest) ([]*models.Order, *common.PaginationResponse, error)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	WithTx(tx *gorm.DB) IOrderRepository
}

type orderRepository struct {
	db *gorm.DB
}

func (or *orderRepository) GetNumbering(ctx context.Context, module string) (*models.Numbering, error) {
	var numbering models.Numbering

	err := or.db.WithContext(ctx).
		Where("module = ?", module).
		First(&numbering).Error
	if err != nil {
		return nil, err
	}

	return &numbering, nil
}

func (or *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return or.db.WithContext(ctx).Create(order).Error
}

func (or *orderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	return or.db.WithContext(ctx).Save(order).Error
}

func (or *orderRepository) UpdateNumbering(ctx context.Context, numbering *models.Numbering) error {
	return or.db.WithContext(ctx).Save(numbering).Error
}

func (or *orderRepository) CreateOrderItem(ctx context.Context, orderItem *models.OrderItem) error {
	return or.db.WithContext(ctx).Create(orderItem).Error
}

func (or *orderRepository) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order

	err := or.db.WithContext(ctx).
		Preload("Items").
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (or *orderRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := or.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (or *orderRepository) WithTx(tx *gorm.DB) IOrderRepository {
	return &orderRepository{
		db: tx,
	}
}

func (or *orderRepository) GetListOrderAdmin(ctx context.Context, pagination *common.PaginationRequest) ([]*models.Order, *common.PaginationResponse, error) {
	var orders []*models.Order
	var totalItems int64

	page := pagination.CurrentPage
	if page < 1 {
		page = 1
	}

	perPage := pagination.PerPage
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	query := or.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("is_deleted = ?", false)

	// Count total items
	err := query.Count(&totalItems).Error
	if err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * perPage

	// Get orders with order items preloaded
	err = or.db.WithContext(ctx).
		Preload("Items").
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&orders).Error
	if err != nil {
		return nil, nil, err
	}

	totalPages := int32(totalItems) / perPage
	if int32(totalItems)%perPage > 0 {
		totalPages++
	}

	paginationResponse := &common.PaginationResponse{
		CurrentPage:    page,
		TotalPageCount: totalPages,
		PerPage:        perPage,
		TotalItemCount: int32(totalItems),
	}

	return orders, paginationResponse, nil
}

func (or *orderRepository) GetListOrder(ctx context.Context, userID uint, pagination *common.PaginationRequest) ([]*models.Order, *common.PaginationResponse, error) {
	var orders []*models.Order
	var totalItems int64

	page := pagination.CurrentPage
	if page < 1 {
		page = 1
	}

	perPage := pagination.PerPage
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	query := or.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("user_id = ?", userID).
		Where("is_deleted = ?", false)

	// Count total items
	err := query.Count(&totalItems).Error
	if err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * perPage

	// Get orders with order items preloaded
	err = or.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&orders).Error
	if err != nil {
		return nil, nil, err
	}

	totalPages := int32(totalItems) / perPage
	if int32(totalItems)%perPage > 0 {
		totalPages++
	}

	paginationResponse := &common.PaginationResponse{
		CurrentPage:    page,
		TotalPageCount: totalPages,
		PerPage:        perPage,
		TotalItemCount: int32(totalItems),
	}

	return orders, paginationResponse, nil
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &orderRepository{
		db: db,
	}
}
