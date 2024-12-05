package services

import (
	"wayra/internal/core/domain/models"
)

type DeliveryService interface {
	Service[models.Delivery]
}
