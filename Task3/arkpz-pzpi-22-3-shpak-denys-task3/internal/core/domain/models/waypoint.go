package models

import "gorm.io/gorm"

type Waypoint struct {
	ID           uint    `gorm:"primaryKey;column:id"`
	Name         string  `gorm:"size:255;not null;column:name"`
	DeviceSerial string  `gorm:"size:255;not null;column:device_serial"`
	Latitude     float64 `gorm:"not null;column:latitude"`
	Longitude    float64 `gorm:"not null;column:longitude"`
	RouteID      uint    `gorm:"not null;column:route_id"`

	Route      Route        `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	SensorData []SensorData `gorm:"foreignKey:WaypointID;constraint:OnDelete:CASCADE;" json:"sensor_data,omitempty"`
}

func (w *Waypoint) LoadRelations(db *gorm.DB) *gorm.DB {
	withRoute := db.Preload("Route").Preload("Route.Company").Preload("Route.Company.Creator")
	withSensorData := withRoute.Preload("SensorData")
	return withSensorData
}
