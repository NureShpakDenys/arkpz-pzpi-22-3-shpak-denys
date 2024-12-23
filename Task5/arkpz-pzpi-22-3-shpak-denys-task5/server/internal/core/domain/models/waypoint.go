package models // import "wayra/internal/core/domain/models"

import "gorm.io/gorm"

// Waypoint is a struct that represents the Waypoint model of the database
type Waypoint struct {
	// ID is the identifier of the waypoint
	// Example: 1
	ID uint `gorm:"primaryKey;column:id"`

	// Name is the name of the waypoint
	// Example: Waypoint 1
	Name string `gorm:"size:255;not null;column:name"`

	// DeviceSerial is the serial number of the device that sent the waypoint
	// Example: 123456789
	DeviceSerial string `gorm:"size:255;not null;column:device_serial"`

	// Latitude is the latitude of the waypoint
	// Example: -12.04318
	Latitude float64 `gorm:"not null;column:latitude"`

	// Longitude is the longitude of the waypoint
	// Example: -77.02824
	Longitude float64 `gorm:"not null;column:longitude"`

	// Altitude is the altitude of the waypoint
	// Example: 0
	RouteID uint `gorm:"not null;column:route_id"`

	// Route is the route to which the waypoint belongs
	Route Route `gorm:"foreignKey:RouteID" json:"route,omitempty"`

	// SensorData is the sensor data of the waypoint
	SensorData []SensorData `gorm:"foreignKey:WaypointID;constraint:OnDelete:CASCADE;" json:"sensor_data,omitempty"`
}

// LoadRelations is an implementation of the LoadRelations interface
func (w *Waypoint) LoadRelations(db *gorm.DB) *gorm.DB {
	withRoute := db.Preload("Route").Preload("Route.Company").Preload("Route.Company.Creator")
	withSensorData := withRoute.Preload("SensorData")
	return withSensorData
}
