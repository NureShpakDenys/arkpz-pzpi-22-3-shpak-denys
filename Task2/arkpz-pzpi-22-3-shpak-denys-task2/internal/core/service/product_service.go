package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type ProductService struct {
	*GenericService[models.Product]
}

func NewProductService(repo port.Repository[models.Product]) *ProductService {
	return &ProductService{
		GenericService: NewGenericService(repo),
	}
}
