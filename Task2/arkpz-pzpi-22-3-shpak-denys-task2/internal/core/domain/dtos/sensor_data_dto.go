package dtos

type SensorDataDTO struct {
	ID          uint    `json:"id,omitempty"`
	Timestamp   string  `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}
