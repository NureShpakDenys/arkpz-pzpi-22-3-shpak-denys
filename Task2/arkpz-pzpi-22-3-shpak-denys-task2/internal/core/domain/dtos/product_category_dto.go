package dtos

type ProductCategoryDTO struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	MinTemperature float64 `json:"min_temperature"`
	MaxTemperature float64 `json:"max_temperature"`
	MinHumidity    float64 `json:"min_humidity"`
	MaxHumidity    float64 `json:"max_humidity"`
	IsPerishable   bool    `json:"is_perishable"`
}
