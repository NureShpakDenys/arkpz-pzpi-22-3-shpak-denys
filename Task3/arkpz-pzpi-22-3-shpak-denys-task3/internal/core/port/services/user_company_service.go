package services

import (
	"wayra/internal/core/domain/models"
)

type UserCompanyService interface {
	Service[models.UserCompany]
	UserBelongsToCompany(userID, companyID uint) bool
}
