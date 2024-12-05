package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type SensorDataService struct {
	*GenericService[models.SensorData]
}

func NewSensorDataService(repo port.Repository[models.SensorData]) *SensorDataService {
	return &SensorDataService{
		GenericService: NewGenericService(repo),
	}
}
