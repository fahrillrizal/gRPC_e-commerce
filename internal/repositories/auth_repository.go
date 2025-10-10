package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/fahrillrizal/ecommerce-grpc/models"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string, updatedBy string) error
	GetRoleByCode(ctx context.Context, code string) (*models.UserRole, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := ar.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		Where("is_deleted = ?", false).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ar *authRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := ar.db.WithContext(ctx).
		Preload("Role").
		Where("is_deleted = ?", false).
		First(&user, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (ar *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	return ar.db.WithContext(ctx).Create(user).Error
}

func (ar *authRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return ar.db.WithContext(ctx).Save(user).Error
}

func (ar *authRepository) UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string, updatedBy string) error {
	return ar.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":   hashedPassword,
			"updated_at": time.Now(),
			"updated_by": updatedBy,
		}).Error
}

func (ar *authRepository) GetRoleByCode(ctx context.Context, code string) (*models.UserRole, error) {
	var role models.UserRole

	err := ar.db.WithContext(ctx).
		Where("code = ?", code).
		Where("is_deleted = ?", false).
		First(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}