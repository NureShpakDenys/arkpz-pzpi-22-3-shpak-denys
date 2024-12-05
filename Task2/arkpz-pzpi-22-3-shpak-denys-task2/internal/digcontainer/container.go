package digcontainer

import (
	"wayra/internal/adapter/config"
	"wayra/internal/adapter/httpserver"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/adapter/repository"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/service"

	"log/slog"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Config
	container.Provide(config.MustLoad)
	container.Provide(slog.Default)

	// Database
	container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return repository.NewGORMDB(cfg.StoragePath)
	})

	// Repositories
	container.Provide(func(db *gorm.DB) port.Repository[models.Company] {
		return repository.NewRepository[models.Company](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.Delivery] {
		return repository.NewRepository[models.Delivery](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.Product] {
		return repository.NewRepository[models.Product](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.ProductCategory] {
		return repository.NewRepository[models.ProductCategory](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.Role] {
		return repository.NewRepository[models.Role](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.Route] {
		return repository.NewRepository[models.Route](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.SensorData] {
		return repository.NewRepository[models.SensorData](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.User] {
		return repository.NewRepository[models.User](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.Waypoint] {
		return repository.NewRepository[models.Waypoint](db)
	})
	container.Provide(func(db *gorm.DB) port.Repository[models.UserCompany] {
		return repository.NewRepository[models.UserCompany](db)
	})

	// Services
	container.Provide(func(repo port.Repository[models.Company]) *service.CompanyService {
		return service.NewCompanyService(repo)
	})
	container.Provide(func(repo port.Repository[models.Delivery]) *service.DeliveryService {
		return service.NewDeliveryService(repo)
	})
	container.Provide(func(repo port.Repository[models.Product]) *service.ProductService {
		return service.NewProductService(repo)
	})
	container.Provide(func(repo port.Repository[models.ProductCategory]) *service.ProductCategoryService {
		return service.NewProductCategoryService(repo)
	})
	container.Provide(func(repo port.Repository[models.Route]) *service.RouteService {
		return service.NewRouteService(repo)
	})
	container.Provide(func(repo port.Repository[models.SensorData]) *service.SensorDataService {
		return service.NewSensorDataService(repo)
	})
	container.Provide(func(repo port.Repository[models.User]) *service.UserService {
		return service.NewUserService(repo)
	})
	container.Provide(func(us *service.UserService, cfg *config.Config) *service.AuthService {
		return service.NewAuthService(us, cfg.AuthConfig.SecretKey, cfg.AuthConfig.TokenExpiry)
	})
	container.Provide(func(repo port.Repository[models.Waypoint]) *service.WaypointService {
		return service.NewWaypointService(repo)
	})
	container.Provide(func(repo port.Repository[models.UserCompany]) *service.UserCompanyService {
		return service.NewUserCompanyService(repo)
	})

	// Handlers
	container.Provide(func(authService *service.AuthService, cfg *config.Config) *handlers.AuthHandler {
		return handlers.NewAuthHandler(authService, cfg.AuthConfig.TokenExpiry)
	})
	container.Provide(func(companyService *service.CompanyService, userCompanyService *service.UserCompanyService) *handlers.CompanyHandler {
		return handlers.NewCompanyHandler(companyService, userCompanyService)
	})
	container.Provide(func(userService *service.UserService) *handlers.UserHandler {
		return handlers.NewUserHandler(userService)
	})
	container.Provide(func(
		routeService *service.RouteService,
		companyService *service.CompanyService,
		userCompanyService *service.UserCompanyService,
	) *handlers.RouteHandler {
		return handlers.NewRoutesHandler(routeService, companyService, userCompanyService)
	})
	container.Provide(func(
		sensorDataService *service.SensorDataService,
		waypointService *service.WaypointService,
		userCompanyService *service.UserCompanyService,
	) *handlers.SensorDataHandler {
		return handlers.NewSensorDataHandler(sensorDataService, waypointService, userCompanyService)
	})
	container.Provide(func(
		waypointService *service.WaypointService,
		routeService *service.RouteService,
		companyService *service.CompanyService,
		userCompanyService *service.UserCompanyService,
	) *handlers.WaypointHandler {
		return handlers.NewWaypointHandler(waypointService, routeService, companyService, userCompanyService)
	})
	container.Provide(func(
		deliveryService *service.DeliveryService,
		companyService *service.CompanyService,
		userCompanyService *service.UserCompanyService,
	) *handlers.DeliveryHandler {
		return handlers.NewDeliveryHandler(deliveryService, companyService, userCompanyService)
	})
	container.Provide(func(
		productService *service.ProductService,
		deliveryService *service.DeliveryService,
		productCategoryService *service.ProductCategoryService,
		companyService *service.CompanyService,
		userCompanyService *service.UserCompanyService,
	) *handlers.ProductHandler {
		return handlers.NewProductHandler(
			productService,
			deliveryService,
			productCategoryService,
			companyService,
			userCompanyService,
		)
	})

	// HTTP Server
	container.Provide(func(
		log *slog.Logger,
		cfg *config.Config,
		authHandler *handlers.AuthHandler,
		companyHandler *handlers.CompanyHandler,
		userHandler *handlers.UserHandler,
		routeHandler *handlers.RouteHandler,
		waypointHandler *handlers.WaypointHandler,
		authService *service.AuthService,
		sensorDataHandler *handlers.SensorDataHandler,
		deliveryHandler *handlers.DeliveryHandler,
		productHandler *handlers.ProductHandler,
	) *gin.Engine {
		return httpserver.NewRouter(
			log,
			cfg,
			authHandler,
			companyHandler,
			userHandler,
			routeHandler,
			waypointHandler,
			authService,
			sensorDataHandler,
			deliveryHandler,
			productHandler,
		)
	})

	return container
}
