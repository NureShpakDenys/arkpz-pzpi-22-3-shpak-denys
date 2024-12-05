package handlers

import (
	"context"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

type WaypointHandler struct {
	waypointService    services.WaypointService
	routeService       services.RouteService
	companyService     services.CompanyService
	userCompanyService services.UserCompanyService
}

func NewWaypointHandler(
	waypointService services.WaypointService,
	routeService services.RouteService,
	companyService services.CompanyService,
	userCompany services.UserCompanyService,
) *WaypointHandler {
	return &WaypointHandler{
		waypointService:    waypointService,
		routeService:       routeService,
		companyService:     companyService,
		userCompanyService: userCompany,
	}
}

type CreateWaypointRequest struct {
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	DeviceSerial string  `json:"device_serial"`
	RouteID      uint    `json:"route_id"`
}

type UpdateWaypointRequest struct {
	Name         string  `json:"name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	DeviceSerial string  `json:"device_serial"`
}

// AddWaypoint godoc
// @Summary      Add a waypoint to a route
// @Description  Adds a new waypoint to the specified route
// @Tags         waypoint
// @Accept       json
// @Produce      json
// @Param        waypoint body CreateWaypointRequest true "Waypoint details"
// @Security     BearerAuth
// @Router       /waypoints [post]
func (h *WaypointHandler) AddWaypoint(c *gin.Context) {
	var waypointRequest CreateWaypointRequest
	if err := c.ShouldBindJSON(&waypointRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	route, err := h.routeService.GetByID(context.Background(), waypointRequest.RouteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	waypoint := &models.Waypoint{
		Name:         waypointRequest.Name,
		DeviceSerial: waypointRequest.DeviceSerial,
		Latitude:     waypointRequest.Latitude,
		Longitude:    waypointRequest.Longitude,
		RouteID:      waypointRequest.RouteID,
	}

	if err := h.waypointService.Create(context.Background(), waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	waypointDTO := &dtos.WaypointDTO{}
	if err = dtoMapper.Map(waypointDTO, waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, waypointDTO)
}

// GetWaypoint godoc
// @Summary      Get waypoint details
// @Description  Retrieves the details of a waypoint
// @Tags         waypoint
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Security     BearerAuth
// @Router       /waypoints/{waypoint_id} [get]
func (h *WaypointHandler) GetWaypoint(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, waypoint.Route.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	waypointDTO := &dtos.WaypointDTO{}
	if err = dtoMapper.Map(waypointDTO, waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, waypointDTO)
}

// UpdateWaypoint godoc
// @Summary      Update waypoint details
// @Description  Updates the details of a waypoint
// @Tags         waypoint
// @Accept       json
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Param        waypoint body UpdateWaypointRequest true "Waypoint details"
// @Security     BearerAuth
// @Router       /waypoints/{waypoint_id} [put]
func (h *WaypointHandler) UpdateWaypoint(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var waypointRequest UpdateWaypointRequest
	if err := c.ShouldBindJSON(&waypointRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	waypoint.Name = waypointRequest.Name
	waypoint.DeviceSerial = waypointRequest.DeviceSerial
	waypoint.Latitude = waypointRequest.Latitude
	waypoint.Longitude = waypointRequest.Longitude
	waypoint.Route = models.Route{}
	waypoint.SensorData = nil

	if err := h.waypointService.Update(context.Background(), waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	waypointDTO := &dtos.WaypointDTO{}
	if err = dtoMapper.Map(waypointDTO, waypoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, waypointDTO)
}

// DeleteWaypoint godoc
// @Summary      Delete waypoint
// @Description  Deletes a waypoint
// @Tags         waypoint
// @Produce      json
// @Param        waypoint_id path int true "Waypoint ID"
// @Security     BearerAuth
// @Router       /waypoints/{waypoint_id} [delete]
func (h *WaypointHandler) DeleteWaypoint(c *gin.Context) {
	waypointID, err := strconv.Atoi(c.Param("waypoint_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid waypoint ID format"})
		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), uint(waypointID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Waypoint not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.waypointService.Delete(context.Background(), uint(waypointID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Waypoint deleted successfully"})
}
