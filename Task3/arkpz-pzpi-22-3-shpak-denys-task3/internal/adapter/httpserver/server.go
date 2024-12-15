package httpserver

import (
	"log/slog"
	"wayra/internal/adapter/config"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/adapter/httpserver/middlewares"
	"wayra/internal/core/port/services"

	_ "wayra/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	log *slog.Logger,
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	companyHandler *handlers.CompanyHandler,
	userHandler *handlers.UserHandler,
	routeHanler *handlers.RouteHandler,
	waypointHandler *handlers.WaypointHandler,
	authService services.AuthService,
	sensorDataHandler *handlers.SensorDataHandler,
	deliveryHandler *handlers.DeliveryHandler,
	productHandler *handlers.ProductHandler,
	adminHandler *handlers.AdminHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterUser)
		auth.POST("/login", authHandler.LoginUser)
	}

	r.Use(middlewares.AuthMiddleware(log, authService))

	r.POST("/auth/logout", authHandler.LogoutUser)

	user := r.Group("/user")
	{
		user.GET("/:id", userHandler.GetUser)
		user.PUT("/:id", userHandler.UpdateUser)
		user.DELETE("/:id", userHandler.DeleteUser)
	}

	r.GET("/users", userHandler.GetUsers)

	company := r.Group("/company")
	{
		company.POST("/", companyHandler.RegisterCompany)
		company.GET("/:company_id", companyHandler.GetCompany)
		company.PUT("/:company_id", companyHandler.UpdateCompany)
		company.DELETE("/:company_id", companyHandler.DeleteCompany)

		company.POST("/:company_id/add-user", companyHandler.AddUserToCompany)
		company.PUT("/:company_id/update-user", companyHandler.UpdateUserInCompany)
		company.DELETE("/:company_id/remove-user", companyHandler.RemoveUserFromCompany)
	}

	deliveries := r.Group("/delivery")
	{
		deliveries.POST("/", deliveryHandler.CreateDelivery)
		deliveries.GET("/:delivery_id", deliveryHandler.GetDelivery)
		deliveries.PUT("/:delivery_id", deliveryHandler.UpdateDelivery)
		deliveries.DELETE("/:delivery_id", deliveryHandler.DeleteDelivery)
	}

	products := r.Group("/products")
	{
		products.POST("/", productHandler.AddProduct)
		products.GET("/:product_id", productHandler.GetProduct)
		products.PUT("/:product_id", productHandler.UpdateProduct)
		products.DELETE("/:product_id", productHandler.DeleteProduct)
	}

	routes := r.Group("/routes")
	{
		routes.POST("/", routeHanler.CreateRoute)
		routes.GET("/:route_id", routeHanler.GetRoute)
		routes.PUT("/:route_id", routeHanler.UpdateRoute)
		routes.DELETE("/:route_id", routeHanler.DeleteRoute)

		routes.GET("/:route_id/weather-alert", routeHanler.GetWeatherAlert)
	}

	analytics := r.Group("/analytics")
	{
		analytics.GET("/:delivery_id/optimal-route", routeHanler.GetOptimalRoute)
		analytics.GET("/:delivery_id/optimal-back-route", routeHanler.GetOptimalBackRoute)
	}

	waypoints := r.Group("/waypoints")
	{
		waypoints.POST("/", waypointHandler.AddWaypoint)
		waypoints.GET("/:waypoint_id", waypointHandler.GetWaypoint)
		waypoints.PUT("/:waypoint_id", waypointHandler.UpdateWaypoint)
		waypoints.DELETE("/:waypoint_id", waypointHandler.DeleteWaypoint)
	}

	sensorData := r.Group("/sensor-data")
	{
		sensorData.POST("/", sensorDataHandler.AddSensorData)
		sensorData.GET("/:sensor_data_id", sensorDataHandler.GetSensorData)
		sensorData.PUT("/:sensor_data_id", sensorDataHandler.UpdateSensorData)
		sensorData.DELETE("/:sensor_data_id", sensorDataHandler.DeleteSensorData)
	}

	admin := r.Group("/admin")
	{
		admin.POST("/backup", adminHandler.BackupDatabase)
		admin.POST("/restore", adminHandler.RestoreDatabase)
	}

	return r
}
