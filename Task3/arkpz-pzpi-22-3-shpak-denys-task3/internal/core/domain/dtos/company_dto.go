package dtos

type CompanyDTO struct {
	ID         uint          `json:"id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	Creator    UserDTO       `json:"creator,omitempty"`
	Users      []UserDTO     `json:"users,omitempty"`
	Routes     []RouteDTO    `json:"routes,omitempty"`
	Deliveries []DeliveryDTO `json:"deliveries,omitempty"`
}
