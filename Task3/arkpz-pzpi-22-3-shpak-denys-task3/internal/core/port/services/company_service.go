package services

import (
	"wayra/internal/core/domain/models"
)

type CompanyService interface {
	Service[models.Company]
}
