package models

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string `gorm:"type:varchar(255);not null" json:"full_name"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	RoleID   *uint  `gorm:"index:idx_user_role_id" json:"role_id"`
	BaseModel
	Role *UserRole `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}
