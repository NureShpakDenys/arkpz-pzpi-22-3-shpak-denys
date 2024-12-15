package models

import "gorm.io/gorm"

type Company struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	Name      string `gorm:"size:255;not null;column:name"`
	Address   string `gorm:"type:text;column:address"`
	CreatorID uint   `gorm:"not null;column:creator_id"`

	Creator    User       `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE;" json:"creator,omitempty"`
	Users      []User     `gorm:"many2many:user_companies;constraint:OnDelete:CASCADE;" json:"users,omitempty"`
	Routes     []Route    `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"routes,omitempty"`
	Deliveries []Delivery `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE;" json:"deliveries,omitempty"`
}

func (c *Company) LoadRelations(db *gorm.DB) *gorm.DB {
	withCreator := db.Preload("Creator")
	withUsers := withCreator.Preload("Users")
	withRoutes := withUsers.Preload("Routes").Preload("Routes.Waypoints").Preload("Routes.Waypoints.SensorData")
	withDeliveries := withRoutes.Preload("Deliveries").Preload("Deliveries.Products").Preload("Deliveries.Products.ProductCategory")
	return withDeliveries
}
