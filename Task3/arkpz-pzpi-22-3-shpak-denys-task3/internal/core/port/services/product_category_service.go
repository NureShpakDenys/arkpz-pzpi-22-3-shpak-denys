package services

import (
	"wayra/internal/core/domain/models"
)

type ProductCategoryService interface {
	Service[models.ProductCategory]
}
