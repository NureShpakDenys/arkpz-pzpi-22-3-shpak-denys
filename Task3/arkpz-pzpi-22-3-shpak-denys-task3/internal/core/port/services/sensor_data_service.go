package services

import (
	"wayra/internal/core/domain/models"
)

type SensorDataService interface {
	Service[models.SensorData]
}
