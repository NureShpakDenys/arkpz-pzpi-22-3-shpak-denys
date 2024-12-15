package dtos // import "wayra/internal/core/domain/dtos"

// WaypointDTO is a data transfer object that represents a Waypoint entity
type WaypointDTO struct {
	// ID is the unique identifier of the Waypoint
	// Example: 1
	ID uint `json:"id,omitempty"`

	// Name is the name of the Waypoint
	// Example: Waypoint 1
	Name string `json:"name"`

	// DeviceSerial is the serial number of the device that sent the Waypoint
	// Example: 123456
	DeviceSerial string `json:"device_serial"`

	// Latitude is the latitude of the Waypoint
	// Example: -12.045
	Latitude float64 `json:"latitude"`

	// Longitude is the longitude of the Waypoint
	// Example: -77.0311
	Longitude float64 `json:"longitude"`

	// Altitude is the altitude of the Waypoint
	SensorData []SensorDataDTO `json:"sensor_data,omitempty"`
}
