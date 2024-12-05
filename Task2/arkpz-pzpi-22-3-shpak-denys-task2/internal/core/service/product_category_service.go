package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type ProductCategoryService struct {
	*GenericService[models.ProductCategory]
}

func NewProductCategoryService(repo port.Repository[models.ProductCategory]) *ProductCategoryService {
	return &ProductCategoryService{
		GenericService: NewGenericService(repo),
	}
}
