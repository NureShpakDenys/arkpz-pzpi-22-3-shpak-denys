package dtos

type WaypointDTO struct {
	ID           uint    `json:"id,omitempty"`
	Name         string  `json:"name"`
	DeviceSerial string  `json:"device_serial"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`

	SensorData []SensorDataDTO `json:"sensor_data,omitempty"`
}
