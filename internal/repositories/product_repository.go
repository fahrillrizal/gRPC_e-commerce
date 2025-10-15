package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/models"
	"github.com/fahrillrizal/ecommerce-grpc/pb/common"
	"gorm.io/gorm"
)

type IProductRepository interface {
	GetProductByID(ctx context.Context, id uint) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, product *models.Product) error
	GetProductsPagination(ctx context.Context, pagination *common.PaginationRequest) ([]models.Product, *common.PaginationResponse, error)
	GetProductsPaginationAdmin(ctx context.Context, pagination *common.PaginationRequest) ([]models.Product, *common.PaginationResponse, error)
	GetHighlightedProducts(ctx context.Context) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func (pr *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return pr.db.WithContext(ctx).Create(product).Error
}

func (pr *productRepository) GetProductByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product

	err := pr.db.WithContext(ctx).
		Where("id = ?", id).
		Where("is_deleted = ?", false).
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	return pr.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Where("is_deleted = ?", false).
		Updates(product).Error
}

func (pr *productRepository) DeleteProduct(ctx context.Context, product *models.Product) error {
	now := time.Now()

	return pr.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Where("is_deleted = ?", false).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": now,
			"deleted_by": product.DeletedBy,
			"updated_at": now,
			"updated_by": product.UpdatedBy,
		}).Error
}

func (pr *productRepository) GetProductsPagination(ctx context.Context, pagination *common.PaginationRequest) ([]models.Product, *common.PaginationResponse, error) {
	var products []models.Product
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

	err := pr.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("is_deleted = ?", false).
		Count(&totalItems).Error

	if err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * perPage

	err = pr.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&products).Error

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

	return products, paginationResponse, nil
}

func (pr *productRepository) GetProductsPaginationAdmin(ctx context.Context, pagination *common.PaginationRequest) ([]models.Product, *common.PaginationResponse, error) {
	var products []models.Product
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

	// ✅ FIX: Tambahkan filter is_deleted = false
	err := pr.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("is_deleted = ?", false).
		Count(&totalItems).Error

	if err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * perPage

	allowedSorts := map[string]bool{
		"name":        true,
		"description": true,
		"price":       true,
		"created_at":  true,
		"updated_at":  true,
	}

	sortClause := "created_at DESC"

	if pagination.Sort != nil {
		sortField := pagination.Sort.Field
		sortDirection := pagination.Sort.Direction

		if sortField != "" && allowedSorts[sortField] {

			if sortDirection == "asc" || sortDirection == "ASC" {
				sortClause = sortField + " ASC"
			} else if sortDirection == "desc" || sortDirection == "DESC" {
				sortClause = sortField + " DESC"
			} else {

				sortClause = sortField + " DESC"
			}
		}
	}

	// ✅ FIX: Tambahkan filter is_deleted = false
	err = pr.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Order(sortClause).
		Limit(int(perPage)).
		Offset(int(offset)).
		Find(&products).Error

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

	return products, paginationResponse, nil
}

func (pr *productRepository) GetHighlightedProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	err := pr.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Limit(10).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func NewProductRepository(db *gorm.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}