package service

import (
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
)

type RoleService struct {
	*GenericService[models.Role]
}

func NewRoleService(repo port.Repository[models.Role]) *RoleService {
	return &RoleService{
		GenericService: NewGenericService(repo),
	}
}
