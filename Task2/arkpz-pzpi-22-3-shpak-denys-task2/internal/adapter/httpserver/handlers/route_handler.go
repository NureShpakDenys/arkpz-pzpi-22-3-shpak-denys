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

type RouteHandler struct {
	routeService       services.RouteService
	companyService     services.CompanyService
	userCompanyService services.UserCompanyService
}

func NewRoutesHandler(
	routeService services.RouteService,
	companyService services.CompanyService,
	userCompanyService services.UserCompanyService,
) *RouteHandler {
	return &RouteHandler{
		routeService:       routeService,
		companyService:     companyService,
		userCompanyService: userCompanyService,
	}
}

const (
	RouteNotStarted = "not_started"
	RouteActive     = "active"
	RouteCompleted  = "completed"
)

type CreateRouteRequest struct {
	Name      string `json:"name" example:"Route 1"`
	Status    string `json:"status" example:"not_started | active | completed"`
	CompanyID uint   `json:"company_id"`
}

type UpdateRouteRequest struct {
	Name   string `json:"name" example:"Route 1"`
	Status string `json:"status" example:"not_started | active | completed"`
}

// CreateRoute godoc
// @Summary      Create a new route
// @Description  Creates a new route with the provided details
// @Tags         route
// @Accept       json
// @Produce      json
// @Param        route body CreateRouteRequest true "Route details"
// @Security     BearerAuth
// @Router       /routes [post]
func (h *RouteHandler) CreateRoute(c *gin.Context) {
	var routeRequest CreateRouteRequest
	if err := c.ShouldBindJSON(&routeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if _, err = h.companyService.GetByID(context.Background(), uint(routeRequest.CompanyID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: routeRequest.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	route := &models.Route{
		Name:      routeRequest.Name,
		CompanyID: uint(routeRequest.CompanyID),
		Status:    routeRequest.Status,
	}

	if err := h.routeService.Create(context.Background(), route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	routeDTO := &dtos.RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeDTO)
}

// GetRoute godoc
// @Summary      Get a route
// @Description  Retrieves a route with the given ID
// @Tags         route
// @Produce      json
// @Param        route_id path int true "Route ID"
// @Security     BearerAuth
// @Router       /routes/{route_id} [get]
func (h *RouteHandler) GetRoute(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Route ID format"})
		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, route.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	routeDTO := &dtos.RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeDTO)
}

// UpdateRoute godoc
// @Summary      Update an existing route
// @Description  Updates an existing route with the given ID
// @Tags         route
// @Accept       json
// @Produce      json
// @Param        route_id path int true "Route ID"
// @Param        route body UpdateRouteRequest true "Updated route details"
// @Security     BearerAuth
// @Router       /routes/{route_id} [put]
func (h *RouteHandler) UpdateRoute(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Route ID format"})
		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found or does not belong to the specified company"})
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

	var routeRequest UpdateRouteRequest
	if err := c.ShouldBindJSON(&routeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if routeRequest.Name != "" {
		route.Name = routeRequest.Name
	}
	if routeRequest.Status != "" {
		route.Status = routeRequest.Status
	}

	route.Company = models.Company{}
	route.Waypoints = nil

	if err := h.routeService.Update(context.Background(), route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	routeDTO := &dtos.RouteDTO{}
	if err = dtoMapper.Map(routeDTO, route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeDTO)
}

// DeleteRoute godoc
// @Summary      Delete a route
// @Description  Deletes a route with the given ID
// @Tags         route
// @Produce      json
// @Param        route_id path int true "Route ID"
// @Security     BearerAuth
// @Router       /routes/{route_id} [delete]
func (h *RouteHandler) DeleteRoute(c *gin.Context) {
	routeID, err := strconv.Atoi(c.Param("route_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID format"})
		return
	}

	route, err := h.routeService.GetByID(context.Background(), uint(routeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found or does not belong to the specified company"})
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

	if err := h.routeService.Delete(context.Background(), uint(routeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route deleted successfully"})
}

// GetOptimalRoute godoc
// @Summary      Get optimal route
// @Description  Retrieves the optimal route for the given route ID
// @Tags         analytics
// @Produce      json
// @Param        company_id path int true "company_id"
// @Security     BearerAuth
// @Router       /analytics/{company_id}/optimal-route [get]
func (h *RouteHandler) GetOptimalRoute(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, uint(companyID)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	// route, err := h.routeService.GetOptimalRoute(context.Background(), uint(companyID))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// routeDTO := &dtos.RouteDTO{}
	// if err = dtoMapper.Map(routeDTO, route); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, map[string]string{"message": "Optimal route"})
}

// GetWeatherAlert godoc
// @Summary      Get weather alert
// @Description  Retrieves the weather alert for the given route ID
// @Tags         analytics
// @Produce      json
// @Param        company_id path int true "company_id"
// @Security     BearerAuth
// @Router       /analytics/{company_id}/weather-alert [get]
func (h *RouteHandler) GetWeatherAlert(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, uint(companyID)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	// weatherAlert, err := h.routeService.GetWeatherAlert(context.Background(), uint(companyID))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// weatherAlertDTO := &dtos.WeatherAlertDTO{}
	// if err = dtoMapper.Map(weatherAlertDTO, weatherAlert); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, map[string]string{"message": "Weather alert"})
}

// GetOptimalBackRoute godoc
// @Summary      Get optimal back route
// @Description  Retrieves the optimal back route for the given route ID
// @Tags         analytics
// @Produce      json
// @Param        company_id path int true "company_id"
// @Security     BearerAuth
// @Router       /analytics/{company_id}/optimal-back-route [get]
func (h *RouteHandler) GetOptimalBackRoute(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, uint(companyID)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	// route, err := h.routeService.GetOptimalBackRoute(context.Background(), uint(companyID))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// routeDTO := &dtos.RouteDTO{}
	// if err = dtoMapper.Map(routeDTO, route); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(http.StatusOK, map[string]string{"message": "Optimal back route"})
}
