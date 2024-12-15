package dtos

type RouteDTO struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name"`

	Waypoints []WaypointDTO `json:"waypoints,omitempty"`
}
