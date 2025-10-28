package models

type OrderItem struct {
	ID           uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uint     `gorm:"not null;index:idx_order_item_order" json:"order_id"`
	Order        *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID    uint     `gorm:"not null;index:idx_order_item_product" json:"product_id"`
	Product      *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName  string   `gorm:"type:varchar(255);not null" json:"product_name"`
	ProductImage string   `gorm:"type:varchar(255)" json:"product_image"`
	ProductPrice float64  `gorm:"type:decimal(15,2);not null" json:"product_price"`
	Quantity     int      `gorm:"not null" json:"quantity"`
	Subtotal     float64  `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	BaseModel
}

func init() {
	RegisterModel(&OrderItem{})
}
