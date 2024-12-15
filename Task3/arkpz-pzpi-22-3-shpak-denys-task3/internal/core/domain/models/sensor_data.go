package models

import (
	"time"

	"gorm.io/gorm"
)

type SensorData struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	Date         time.Time `gorm:"type:timestamp;not null;column:date"`
	Temperature  float64   `gorm:"not null;column:temperature"`
	Humidity     float64   `gorm:"not null;column:humidity"`
	WindSpeed    float64   `gorm:"not null;column:wind_speed"`
	MeanPressure float64   `gorm:"not null;column:mean_pressure"`
	WaypointID   uint      `gorm:"not null;column:waypoint_id"`

	Waypoint Waypoint `gorm:"foreignKey:WaypointID" json:"device,omitempty"`
}

func (s *SensorData) LoadRelations(db *gorm.DB) *gorm.DB {
	//_ = db.Preload("Waypoint")
	return db
}
