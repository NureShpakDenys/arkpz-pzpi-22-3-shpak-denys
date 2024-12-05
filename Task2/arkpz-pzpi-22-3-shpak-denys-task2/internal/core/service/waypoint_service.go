package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type WaypointService struct {
	*GenericService[models.Waypoint]
}

func NewWaypointService(repo port.Repository[models.Waypoint]) *WaypointService {
	return &WaypointService{
		GenericService: NewGenericService(repo),
	}
}
