package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(255)" json:"created_by"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(255)" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`
	DeletedBy *string        `gorm:"type:varchar(255)" json:"deleted_by,omitempty"`
	IsDeleted bool           `gorm:"type:boolean;not null;default:false" json:"is_deleted"`
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	if b.DeletedAt.Valid && !b.IsDeleted {
		b.IsDeleted = true
	}
	return nil
}
