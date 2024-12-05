package services

import (
	"context"
	"wayra/internal/core/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	ValidateToken(token string) (*jwt.RegisteredClaims, error)
}
