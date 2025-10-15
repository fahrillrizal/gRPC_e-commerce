package models

type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Price       float64 `gorm:"type:decimal(10,2);not null;uniqueIndex" json:"price"`
	Description string  `gorm:"type:text" json:"description"`
	ImageURL    string  `gorm:"type:varchar(255)" json:"image_url"`
	BaseModel
}