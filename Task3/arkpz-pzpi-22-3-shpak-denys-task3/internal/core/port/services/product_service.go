package services

import (
	"wayra/internal/core/domain/models"
)

type ProductService interface {
	Service[models.Product]
}
