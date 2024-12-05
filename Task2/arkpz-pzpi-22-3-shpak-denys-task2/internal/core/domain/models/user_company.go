package models

import "gorm.io/gorm"

type UserCompany struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	UserID    uint   `gorm:"not null;column:user_id"`
	CompanyID uint   `gorm:"not null;column:company_id"`
	Role      string `gorm:"size:50;default:'user';column:role"`

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Company Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company,omitempty"`
}

func (u *UserCompany) LoadRelations(db *gorm.DB) *gorm.DB {
	withUser := db.Preload("User")
	withCompany := withUser.Preload("Company")
	return withCompany
}
