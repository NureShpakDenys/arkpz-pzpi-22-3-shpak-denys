package dtos

import "time"

type DeliveryDTO struct {
	ID        uint         `json:"id"`
	Status    string       `json:"status"`
	Date      time.Time    `json:"date"`
	Duration  string       `json:"duration"`
	CompanyID uint         `json:"company_id"`
	RouteID   uint         `json:"route_id"`
	Products  []ProductDTO `json:"products,omitempty"`
}
