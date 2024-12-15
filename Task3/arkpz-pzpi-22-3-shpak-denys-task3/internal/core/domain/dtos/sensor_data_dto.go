package dtos

import "time"

type SensorDataDTO struct {
	ID           uint      `json:"id,omitempty"`
	Date         time.Time `json:"date"`
	Temperature  float64   `json:"temperature"`
	Humidity     float64   `json:"humidity"`
	WindSpeed    float64   `json:"wind_speed"`
	MeanPressure float64   `json:"mean_pressure"`
}
