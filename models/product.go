package models

type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`
	Price       float64 `gorm:"type:decimal(15,2);not null;" json:"price"`
	Description string  `gorm:"type:text" json:"description"`
	ImageURL    string  `gorm:"type:varchar(255)" json:"image_url"`
	BaseModel
}