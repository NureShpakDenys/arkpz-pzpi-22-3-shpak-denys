package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

type SensorDataHandler struct {
	sensorDataService  services.SensorDataService
	waypointService    services.WaypointService
	userCompanyService services.UserCompanyService
}

func NewSensorDataHandler(
	sensorDataService services.SensorDataService,
	waypointService services.WaypointService,
	userCompanyService services.UserCompanyService,
) *SensorDataHandler {
	return &SensorDataHandler{
		sensorDataService:  sensorDataService,
		waypointService:    waypointService,
		userCompanyService: userCompanyService,
	}
}

type CreateSensorDataRequest struct {
	Date        string  `json:"date" example:"2021-09-01T12:00:00Z"`
	Temperature float64 `json:"temperature" example:"25.5"`
	Humidity    float64 `json:"humidity" example:"50.0"`
	WaypointID  uint    `json:"waypoint_id"`
}

type UpdateSensorDataRequest struct {
	Date        string  `json:"date" example:"2021-09-01T12:00:00Z"`
	Temperature float64 `json:"temperature" example:"25.5"`
	Humidity    float64 `json:"humidity" example:"50.0"`
}

// AddSensorData godoc
// @Summary      Add sensor data to a SensorData
// @Description  Adds new sensor data to the specified SensorData
// @Tags         sensor
// @Accept       json
// @Produce      json
// @Param        sensor_data body CreateSensorDataRequest true "Sensor data details"
// @Security     BearerAuth
// @Router       /sensor-data [post]
func (h *SensorDataHandler) AddSensorData(c *gin.Context) {
	var sensorDataRequest CreateSensorDataRequest
	if err := c.ShouldBindJSON(&sensorDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	waypoint, err := h.waypointService.GetByID(context.Background(), sensorDataRequest.WaypointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	var date time.Time
	if sensorDataRequest.Date != "" {
		date, err = time.Parse(time.RFC3339, sensorDataRequest.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	sensorData := &models.SensorData{
		Date:        date,
		Temperature: sensorDataRequest.Temperature,
		Humidity:    sensorDataRequest.Humidity,
		WaypointID:  sensorDataRequest.WaypointID,
	}

	if err := h.sensorDataService.Create(context.Background(), sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sensorDataDTO := &dtos.SensorDataDTO{}
	if err = dtoMapper.Map(sensorDataDTO, sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorDataDTO)
}

// GetSensorData godoc
// @Summary      Get sensor data by ID
// @Description  Retrieves sensor data with the given ID
// @Tags         sensor
// @Produce      json
// @Param        sensor_data_id path int true "Sensor Data ID"
// @Security     BearerAuth
// @Router       /sensor-data/{sensor_data_id} [get]
func (h *SensorDataHandler) GetSensorData(c *gin.Context) {
	sensorDataID, err := strconv.Atoi(c.Param("sensor_data_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sensor Data ID format"})
		return
	}

	sensorData, err := h.sensorDataService.GetByID(context.Background(), uint(sensorDataID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, sensorData.Waypoint.Route.CompanyID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this sensorData"})
		return
	}

	sensorDataDTO := &dtos.SensorDataDTO{}
	if err = dtoMapper.Map(sensorDataDTO, sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorDataDTO)
}

func (h *SensorDataHandler) UpdateSensorData(c *gin.Context) {
	sensorDataID, err := strconv.Atoi(c.Param("sensor_data_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sensor Data ID format"})
		return
	}

	sensorData, err := h.sensorDataService.GetByID(context.Background(), uint(sensorDataID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: sensorData.Waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	var sensorDataRequest UpdateSensorDataRequest
	if err := c.ShouldBindJSON(&sensorDataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var date time.Time
	if sensorDataRequest.Date != "" {
		date, err = time.Parse(time.RFC3339, sensorDataRequest.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
	}

	sensorData.Date = date
	sensorData.Temperature = sensorDataRequest.Temperature
	sensorData.Humidity = sensorDataRequest.Humidity
	sensorData.Waypoint = models.Waypoint{}

	if err := h.sensorDataService.Update(context.Background(), sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sensorDataDTO := &dtos.SensorDataDTO{}
	if err = dtoMapper.Map(sensorDataDTO, sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensorDataDTO)
}

func (h *SensorDataHandler) DeleteSensorData(c *gin.Context) {
	sensorDataID, err := strconv.Atoi(c.Param("sensor_data_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Sensor Data ID format"})
		return
	}

	sensorData, err := h.sensorDataService.GetByID(context.Background(), uint(sensorDataID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor Data not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: sensorData.Waypoint.Route.CompanyID,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.sensorDataService.Delete(context.Background(), uint(sensorDataID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor Data deleted successfully"})
}
