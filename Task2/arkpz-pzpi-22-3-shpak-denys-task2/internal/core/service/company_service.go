package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type CompanyService struct {
	*GenericService[models.Company]
}

func NewCompanyService(repo port.Repository[models.Company]) *CompanyService {
	return &CompanyService{
		GenericService: NewGenericService(repo),
	}
}
