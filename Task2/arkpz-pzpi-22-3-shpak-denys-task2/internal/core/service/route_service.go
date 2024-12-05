package service

import (
	"context"
	"errors"
	"math"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type RouteService struct {
	*GenericService[models.Route]
}

func NewRouteService(repo port.Repository[models.Route]) *RouteService {
	return &RouteService{
		GenericService: NewGenericService(repo),
	}
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (s *RouteService) GetOptimalRoute(ctx context.Context, companyID uint) (models.Route, error) {
	routes, err := s.Where(ctx, &models.Route{CompanyID: companyID})
	if err != nil || len(routes) == 0 {
		return models.Route{}, errors.New("no routes found for the company")
	}

	var optimalRoute models.Route
	minDistance := math.MaxFloat64

	for _, route := range routes {
		totalDistance := 0.0
		previousWaypoint := models.Waypoint{}
		validRoute := true

		for i, waypoint := range route.Waypoints {
			if i > 0 {
				distance := haversineDistance(previousWaypoint.Latitude, previousWaypoint.Longitude, waypoint.Latitude, waypoint.Longitude)
				totalDistance += distance
			}
			previousWaypoint = waypoint

			for _, sensorData := range waypoint.SensorData {
				if sensorData.Temperature < -10 || sensorData.Temperature > 35 || sensorData.Humidity > 90 {
					validRoute = false
					break
				}
			}
			if !validRoute {
				break
			}
		}

		if validRoute && totalDistance < minDistance {
			minDistance = totalDistance
			optimalRoute = route
		}
	}

	if optimalRoute.ID == 0 {
		return models.Route{}, errors.New("no optimal route could be determined")
	}
	return optimalRoute, nil
}
