package services

import (
	"wayra/internal/core/domain/models"
)

type RoleService interface {
	Service[models.Role]
}
