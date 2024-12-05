package services

import (
	"wayra/internal/core/domain/models"
)

type UserService interface {
	Service[models.User]
}
