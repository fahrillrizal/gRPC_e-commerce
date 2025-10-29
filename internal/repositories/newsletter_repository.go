package repositories

import (
	"context"

	"github.com/fahrillrizal/ecommerce-grpc/models"
	"gorm.io/gorm"
)

type INewsletterRepository interface {
	GetNewsletterByEmail(ctx context.Context, email string) (*models.Newsletter, error)
	CreateNewsletter(ctx context.Context, newsletter *models.Newsletter) error
}

type newsletterRepository struct {
	db *gorm.DB
}

func (nr *newsletterRepository) GetNewsletterByEmail(ctx context.Context, email string) (*models.Newsletter, error) {
	var newsletter models.Newsletter
	if err := nr.db.WithContext(ctx).Where("email = ?", email).First(&newsletter).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &newsletter, nil
}

func (nr *newsletterRepository) CreateNewsletter(ctx context.Context, newsletter *models.Newsletter) error {
	if err := nr.db.WithContext(ctx).Create(newsletter).Error; err != nil {
		return err
	}
	return nil
}

func NewNewsletterRepository(db *gorm.DB) INewsletterRepository {
	return &newsletterRepository{
		db: db,
	}
}
