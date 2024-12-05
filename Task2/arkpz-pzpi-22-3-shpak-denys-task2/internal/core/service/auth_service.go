package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService services.UserService
	secretKey   string
	tokenExpiry time.Duration
}

func NewAuthService(userService services.UserService, secretKey string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userService: userService,
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
	}
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userService.Create(ctx, user)
}

func (s *AuthService) LoginUser(ctx context.Context, username, password string) (string, error) {
	users, err := s.userService.Where(ctx, &models.User{Name: username})
	if err != nil || len(users) == 0 {
		return "", errors.New("invalid username or password")
	}

	user := users[0]
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ValidateToken(token string) (*jwt.RegisteredClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims.RegisteredClaims, nil
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(int(userID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}
