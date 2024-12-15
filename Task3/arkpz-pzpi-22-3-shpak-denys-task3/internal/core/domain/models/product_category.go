package models

type ProductCategory struct {
	ID             uint    `gorm:"primaryKey;column:id"`
	Name           string  `gorm:"size:255;not null;column:name"`
	Description    string  `gorm:"type:text;column:description"`
	MinTemperature float64 `gorm:"not null;column:min_temperature"`
	MaxTemperature float64 `gorm:"not null;column:max_temperature"`
	MinHumidity    float64 `gorm:"not null;column:min_humidity"`
	MaxHumidity    float64 `gorm:"not null;column:max_humidity"`
	IsPerishable   bool    `gorm:"not null;column:is_perishable"`
}
