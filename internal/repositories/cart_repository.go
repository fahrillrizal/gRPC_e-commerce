package repositories

import (
	"context"
	"errors"

	"github.com/fahrillrizal/ecommerce-grpc/models"
	"gorm.io/gorm"
)

type ICartRepository interface {
	GetCartByProductUserID(ctx context.Context, productId uint, userId uint) (*models.Cart, error)
	CreateNewCart(ctx context.Context, cart *models.Cart) error
	UpdateCart(ctx context.Context, cart *models.Cart) error
	GetListCart(ctx context.Context, userId uint) ([]*models.Cart, error)
	GetCartById(ctx context.Context, cartId uint) (*models.Cart, error)
	DeleteCart(ctx context.Context, cartId uint, deletedBy string) error
}

func (cr *cartRepository) GetCartByProductUserID(ctx context.Context, productId uint, userId uint) (*models.Cart, error) {
	var cart models.Cart

	err := cr.db.WithContext(ctx).
		Preload("Product").
		Preload("User").
		Where("product_id = ? AND user_id = ?", productId, userId).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (cr *cartRepository) CreateNewCart(ctx context.Context, cart *models.Cart) error {
	err := cr.db.WithContext(ctx).Create(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) UpdateCart(ctx context.Context, cart *models.Cart) error {
	err := cr.db.WithContext(ctx).Save(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) GetListCart(ctx context.Context, userId uint) ([]*models.Cart, error) {
	var carts []*models.Cart

	err := cr.db.WithContext(ctx).
		Preload("Product").
		Preload("User").
		Where("user_id = ?", userId).
		Where("is_deleted = ?", false).
		Find(&carts).Error

	if err != nil {
		return nil, err
	}

	return carts, nil
}

func (cr *cartRepository) GetCartById(ctx context.Context, cartId uint) (*models.Cart, error) {
	var cart models.Cart
	err := cr.db.WithContext(ctx).
		Preload("Product").
		Preload("User").
		Where("id = ?", cartId).
		Where("is_deleted = ?", false).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (cr *cartRepository) DeleteCart(ctx context.Context, cartId uint, deletedBy string) error {
	return cr.db.WithContext(ctx).
		Model(&models.Cart{}).
		Where("id = ?", cartId).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_by": deletedBy,
		}).Error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) ICartRepository {
	return &cartRepository{
		db: db,
	}
}
