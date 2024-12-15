package models

import (
	"time"

	"gorm.io/gorm"
)

type Delivery struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Status    string    `gorm:"size:50;column:status"`
	Date      time.Time `gorm:"type:timestamp;not null;column:date"`
	Duration  string    `gorm:"type:interval;not null;column:duration"`
	CompanyID uint      `gorm:"not null;column:company_id"`
	RouteID   uint      `gorm:"not null;column:route_id"`

	Route    Route     `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE;" json:"route,omitempty"`
	Company  Company   `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"company,omitempty"`
	Products []Product `gorm:"foreignKey:DeliveryID;constraint:OnDelete:CASCADE;" json:"products,omitempty"`
}

func (d *Delivery) LoadRelations(db *gorm.DB) *gorm.DB {
	withCompany := db.Preload("Company").Preload("Company.Creator").Preload("Company.Users")
	withRoute := withCompany.Preload("Route").Preload("Route.Waypoints")
	withProducts := withRoute.Preload("Products").Preload("Products.ProductCategory")
	return withProducts
}
