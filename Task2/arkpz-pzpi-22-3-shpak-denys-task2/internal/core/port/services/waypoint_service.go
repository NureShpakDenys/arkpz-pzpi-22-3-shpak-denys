package services

import (
	"wayra/internal/core/domain/models"
)

type WaypointService interface {
	Service[models.Waypoint]
}
