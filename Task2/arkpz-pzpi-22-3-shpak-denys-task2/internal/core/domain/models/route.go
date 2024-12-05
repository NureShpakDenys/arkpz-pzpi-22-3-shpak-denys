package models

import "gorm.io/gorm"

type Route struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	Name      string `gorm:"size:255;not null;column:name"`
	CompanyID uint   `gorm:"not null;column:company_id"`
	Status    string `gorm:"size:50;default:'not_started';column:status"`

	Company   Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Waypoints []Waypoint `gorm:"foreignKey:RouteID;constraint:OnDelete:CASCADE;" json:"waypoints,omitempty"`
}

func (r *Route) LoadRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("Company").Preload("Company.Creator").Preload("Waypoints").Preload("Waypoints.SensorData")
}
