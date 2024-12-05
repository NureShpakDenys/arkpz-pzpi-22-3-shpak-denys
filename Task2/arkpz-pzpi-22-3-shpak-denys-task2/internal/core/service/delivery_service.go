package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type DeliveryService struct {
	*GenericService[models.Delivery]
}

func NewDeliveryService(repo port.Repository[models.Delivery]) *DeliveryService {
	return &DeliveryService{
		GenericService: NewGenericService(repo),
	}
}
