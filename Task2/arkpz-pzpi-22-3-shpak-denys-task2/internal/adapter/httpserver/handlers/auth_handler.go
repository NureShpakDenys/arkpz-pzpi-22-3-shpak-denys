package handlers

import (
	"context"
	"net/http"
	"time"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
)

const (
	AdminRole = iota
	UserRole
)

type AuthHandler struct {
	authService services.AuthService
	tokenExpiry time.Duration
}

func NewAuthHandler(authService services.AuthService, tokenExpiry time.Duration) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		tokenExpiry: tokenExpiry,
	}
}

type AuthCredentials struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"password123"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Registers a new user with the provided details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body AuthCredentials true "User details"
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var request AuthCredentials
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Username == "" || len(request.Username) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 4 characters long"})
		return
	}

	if request.Password == "" || len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	user := models.User{
		Name:     request.Username,
		Password: request.Password,
		RoleID:   UserRole,
	}

	if err := h.authService.RegisterUser(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticates a user and returns a token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body AuthCredentials true "User credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var credentials AuthCredentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.LoginUser(context.Background(), credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.SetCookie("token", token, int(h.tokenExpiry.Seconds()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// LogoutUser godoc
// @Summary      Logout user
// @Description  Logs out a user by invalidating their token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
