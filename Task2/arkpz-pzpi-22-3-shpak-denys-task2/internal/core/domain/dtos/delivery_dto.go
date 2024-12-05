package dtos

type DeliveryDTO struct {
	ID        uint         `json:"id"`
	Status    string       `json:"status"`
	Date      string       `json:"date"`
	CompanyID uint         `json:"company_id"`
	Products  []ProductDTO `json:"products,omitempty"`
}
