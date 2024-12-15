package service

import (
	"context"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"

	"gorm.io/gorm"
)

type UserCompanyService struct {
	*GenericService[models.UserCompany]
}

func NewUserCompanyService(repo port.Repository[models.UserCompany]) *UserCompanyService {
	return &UserCompanyService{
		GenericService: NewGenericService(repo),
	}
}

func (s *UserCompanyService) UserBelongsToCompany(userID, companyID uint) bool {
	userCompany, err := s.getByUserAndCompany(context.Background(), userID, companyID)
	if err != nil || userCompany == nil {
		return false
	}
	return true
}

func (s *UserCompanyService) getByUserAndCompany(ctx context.Context, userID, companyID uint) (*models.UserCompany, error) {
	userCompanies, err := s.Where(ctx, &models.UserCompany{
		UserID:    userID,
		CompanyID: companyID,
	})
	if err != nil {
		return nil, err
	}

	if len(userCompanies) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &userCompanies[0], nil
}
