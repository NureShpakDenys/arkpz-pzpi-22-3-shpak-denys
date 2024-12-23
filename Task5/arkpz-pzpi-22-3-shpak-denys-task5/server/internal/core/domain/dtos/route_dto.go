package dtos // import "wayra/internal/core/domain/dtos"

// RouteDTO is a DTO that represents a route
type RouteDTO struct {
	// ID is the unique identifier of the route
	// Example: 1
	ID uint `json:"id,omitempty"`

	// Name is the name of the route
	// Example: Route 1
	Name string `json:"name"`

	// Waypoints is a list of waypoints that the route has
	Waypoints []WaypointDTO `json:"waypoints,omitempty"`
}
