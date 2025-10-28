package models

type Numbering struct {
	Module string `gorm:"primaryKey;type:varchar(255)" json:"module"`
	Number int    `gorm:"type:int;not null" json:"number"`
}

func init() {
	RegisterModel(&Numbering{})
}