package models

import "gorm.io/gorm"

type User struct {
	ID       uint   `gorm:"primaryKey;column:id"`
	Name     string `gorm:"size:255;not null;column:name" json:"name"`
	Password string `gorm:"size:255;not null;column:password"`
	RoleID   uint   `gorm:"not null;column:role_id"`

	Role      Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Companies []Company `gorm:"many2many:user_companies;constraint:OnDelete:CASCADE;" json:"companies,omitempty"`
}

func (u *User) LoadRelations(db *gorm.DB) *gorm.DB {
	withRole := db.Preload("Role")
	withCompanies := withRole.Preload("Companies").Preload("Companies.Creator").Preload("Companies.Users")
	return withCompanies
}
