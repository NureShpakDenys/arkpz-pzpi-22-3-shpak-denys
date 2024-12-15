package repository

import (
	"wayra/internal/core/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGORMDB(connectionString string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connectionString), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Company{},
		&models.Delivery{},
		&models.Role{},
		&models.Route{},
		&models.SensorData{},
		&models.User{},
		&models.Waypoint{},
		&models.UserCompany{},
		&models.Product{},
		&models.ProductCategory{},
	)
}
