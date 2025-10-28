package models

import "time"

type Order struct {
	ID                   uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Number               string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"number"`
	UserID               uint         `gorm:"not null;index:idx_order_user" json:"user_id"`
	User                 *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderStatusCode      string       `gorm:"type:varchar(50);not null;index:idx_order_status" json:"order_status_code"`
	OrderStatus          *OrderStatus `gorm:"foreignKey:OrderStatusCode;references:Code" json:"order_status,omitempty"`
	UserFullName         string       `gorm:"type:varchar(255);not null" json:"user_full_name"`
	Address              string       `gorm:"type:text;not null" json:"address"`
	PhoneNumber          string       `gorm:"type:varchar(20);not null" json:"phone_number"`
	Notes                string       `gorm:"type:text" json:"notes"`
	Total                float64      `gorm:"type:decimal(15,2);not null" json:"total"`
	ExpiredAt            *time.Time   `gorm:"type:timestamptz;index:idx_order_expired" json:"expired_at,omitempty"`
	XenditInvoiceID      string       `gorm:"type:varchar(255);uniqueIndex" json:"xendit_invoice_id"`
	XenditInvoiceUrl     string       `gorm:"type:varchar(255)" json:"xendit_invoice_url"`
	XenditPaidAt         *time.Time   `gorm:"type:timestamptz" json:"xendit_paid_at,omitempty"`
	XenditPaymentMethod  string       `gorm:"type:varchar(255)" json:"xendit_payment_method,omitempty"`
	XenditPaymentChannel string       `gorm:"type:varchar(255)" json:"xendit_payment_channel,omitempty"`
	Items                []*OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	BaseModel
}

func init() {
	RegisterModel(&Order{})
}
