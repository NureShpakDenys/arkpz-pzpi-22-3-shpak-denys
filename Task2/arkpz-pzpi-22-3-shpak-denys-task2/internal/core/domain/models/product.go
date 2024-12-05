package models

import "gorm.io/gorm"

type Product struct {
	ID                uint    `gorm:"primaryKey;column:id"`
	Name              string  `gorm:"size:255;not null;column:name"`
	Weight            float64 `gorm:"not null;column:weight"`
	ProductCategoryID uint    `gorm:"not null;column:product_category_id"`
	DeliveryID        uint    `gorm:"not null;column:delivery_id"`

	ProductCategory ProductCategory `gorm:"foreignKey:ProductCategoryID" json:"product_category,omitempty"`
	Delivery        Delivery        `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE;" json:"delivery,omitempty"`
}

func (d *Product) LoadRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("ProductCategory").Preload("Delivery")
}
