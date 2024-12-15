package models

type WeatherAlert struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details string `json:"details"`
}
