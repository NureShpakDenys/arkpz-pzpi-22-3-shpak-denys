package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type UserService struct {
	*GenericService[models.User]
}

func NewUserService(repo port.Repository[models.User]) *UserService {
	return &UserService{
		GenericService: NewGenericService(repo),
	}
}
