package service // import "wayra/internal/core/service"

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

// CompanyService is a struct that defines the CompanyService
type CompanyService struct {
	*GenericService[models.Company] // Embedding the GenericService
}

// NewCompanyService is a function that returns a new CompanyService
// port: port.Repository[models.Company] - The repository that will be used by the service
// returns: *CompanyService - The service that will be used to interact with the repository
func NewCompanyService(repo port.Repository[models.Company]) *CompanyService {
	return &CompanyService{
		GenericService: NewGenericService(repo),
	}
}
