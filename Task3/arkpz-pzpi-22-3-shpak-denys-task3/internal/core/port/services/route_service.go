package services

import (
	"context"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/domain/utils/analysis"
)

type RouteService interface {
	Service[models.Route]
	GetOptimalRoute(
		ctx context.Context,
		delivery *models.Delivery,
		includeWeight bool,
		considerPerishable bool,
	) (string, *analysis.PredictData, []float64, models.Route, error)
	GetWeatherAlert(ctx context.Context, route models.Route) ([]models.WeatherAlert, error)
}
