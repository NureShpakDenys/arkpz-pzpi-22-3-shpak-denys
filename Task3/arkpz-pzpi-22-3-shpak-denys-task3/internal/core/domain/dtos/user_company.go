package dtos

type UserCompanyDTO struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	CompanyID uint   `json:"company_id"`
	Role      string `json:"role"`

	User    UserDTO    `json:"user,omitempty"`
	Company CompanyDTO `json:"company,omitempty"`
}
