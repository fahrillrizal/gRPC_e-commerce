package models

type UserRole struct {
	ID	uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"type:varchar(255);not null" json:"name"`
    Code  string `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
    BaseModel
    
    Users []User `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL" json:"users,omitempty"`
}

func init() {
	RegisterModel(&UserRole{})
}