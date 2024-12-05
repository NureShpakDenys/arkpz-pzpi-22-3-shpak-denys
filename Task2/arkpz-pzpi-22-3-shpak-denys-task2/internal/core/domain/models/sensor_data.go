package models

import "gorm.io/gorm"

type SensorData struct {
	ID          uint    `gorm:"primaryKey;column:id"`
	Timestamp   string  `gorm:"type:timestamp;not null;column:timestamp"`
	Temperature float64 `gorm:"not null;column:temperature"`
	Humidity    float64 `gorm:"not null;column:humidity"`
	WaypointID  uint    `gorm:"not null;column:waypoint_id"`

	Waypoint Waypoint `gorm:"foreignKey:WaypointID" json:"device,omitempty"`
}

func (s *SensorData) LoadRelations(db *gorm.DB) *gorm.DB {
	withWaypoint := db.Preload("Waypoint")
	return withWaypoint
}
