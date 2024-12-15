package models

type Role struct {
	ID   uint   `gorm:"primaryKey;column:id"`
	Name string `gorm:"size:255;not null;column:name" json:"name"`

	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}
