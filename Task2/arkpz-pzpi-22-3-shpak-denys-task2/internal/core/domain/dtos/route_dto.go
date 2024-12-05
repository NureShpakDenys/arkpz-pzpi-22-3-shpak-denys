package dtos

type RouteDTO struct {
	ID     uint   `json:"id,omitempty"`
	Name   string `json:"name"`
	Status string `json:"status"`

	Waypoints []WaypointDTO `json:"waypoints,omitempty"`
}
