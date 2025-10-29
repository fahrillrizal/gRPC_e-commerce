package models

type Newsletter struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"type:varchar(255);not null" json:"full_name"`
	Email    string `gorm:"type:varchar(255);not null;" json:"email"`
	BaseModel
}

func init() {
	RegisterModel(&Newsletter{})
}