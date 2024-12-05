package models

import "gorm.io/gorm"

type Delivery struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	Status    string `gorm:"size:50;column:status"`
	Date      string `gorm:"type:timestamp;not null;column:date"`
	CompanyID uint   `gorm:"not null;column:company_id"`

	Company  Company   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company,omitempty"`
	Products []Product `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE;" json:"products,omitempty"`
}

func (d *Delivery) LoadRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("Company").Preload("Company.Creator").Preload("Company.Users").Preload("Products").Preload("Products.ProductCategory")
}
