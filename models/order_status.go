package models

type OrderStatus struct {
	ID          uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string   `gorm:"type:varchar(100);not null" json:"name"`
	Code        string   `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string   `gorm:"type:text" json:"description"`
	Orders      []*Order `gorm:"foreignKey:OrderStatusCode;references:Code" json:"orders,omitempty"`
	BaseModel
}

func init() {
	RegisterModel(&OrderStatus{})
}
