package services

import (
	"context"
	"wayra/internal/core/domain/models"
)

type RouteService interface {
	Service[models.Route]
	GetOptimalRoute(ctx context.Context, companyID uint) (models.Route, error)
}
