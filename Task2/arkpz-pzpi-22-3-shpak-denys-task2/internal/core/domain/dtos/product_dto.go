package dtos

type ProductDTO struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Weight float64 `json:"weight"`

	ProductCategory ProductCategoryDTO `json:"product_category,omitempty"`
}
