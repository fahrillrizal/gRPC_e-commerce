package models

type Cart struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint    `gorm:"not null;index:idx_cart_product_user,unique" json:"product_id"`
	UserID      uint    `gorm:"not null;index:idx_cart_product_user,unique" json:"user_id"`
	Quantity    int    `gorm:"not null" json:"quantity"`
	BaseModel
	Product	 *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	User      *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func init() {
	RegisterModel(&Cart{})
}