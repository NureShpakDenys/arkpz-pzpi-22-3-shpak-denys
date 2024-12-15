package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// CompanyHandler is a struct that handles company-related HTTP requests
type CompanyHandler struct {
	companyService     services.CompanyService     // Company service is an interface for comapny business logic
	userCompanyService services.UserCompanyService // UserCompany service is an interface for user-company business logic
}

// NewCompanyHandler creates a new CompanyHandler with the provided services
// companyService: Service for company operations
// userCompanyService: Service for user-company operations
// Returns: A new CompanyHandler
func NewCompanyHandler(companyService services.CompanyService, userCompanyService services.UserCompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService:     companyService,
		userCompanyService: userCompanyService,
	}
}

// CompanyRequest is the request for the RegisterCompany endpoint
type CompanyRequest struct {
	// Name is the name of the company
	// Example: Wayra
	Name string `json:"name"`

	// Address is the address of the company
	// Example: 123 Main St
	Address string `json:"address"`
}

// AddUserToCompanyRequest is the request for the AddUserToCompany endpoint
type AddUserToCompanyRequest struct {
	// UserID is the ID of the user to add
	// Example: 1
	UserID uint `json:"userID"`

	// Role is the role of the user in the company
	Role Role `json:"role" example:"user | admin | manager"`
}

// Role is the role of a user in a company
type UpdateUserInCompanyRequest struct {
	// UserID is the ID of the user to update
	// Example: 1
	UserID uint `json:"userID"`

	// Role is the role of the user in the company
	Role string `json:"role" example:"user | admin | manager"`
}

// RemoveUserFromCompanyRequest is the request for the RemoveUserFromCompany endpoint
type RemoveUserFromCompanyRequest struct {
	// UserID is the ID of the user to remove
	// Example: 1
	UserID uint `json:"userID"`
}

// RegisterCompany godoc
// @Summary      Register a new company
// @Description  Registers a new company with the provided details
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company body CompanyRequest true "Company details"
// @Security     BearerAuth
// @Router       /company [post]
func (h *CompanyHandler) RegisterCompany(c *gin.Context) {
	var companyRequest CompanyRequest
	if err := c.ShouldBindJSON(&companyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if companyRequest.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if companyRequest.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	company := &models.Company{
		Name:      companyRequest.Name,
		Address:   companyRequest.Address,
		CreatorID: *userID,
	}
	if err := h.companyService.Create(context.Background(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = h.userCompanyService.Create(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: company.ID,
		Role:      string(RoleAdmin),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	companyDTO := &dtos.CompanyDTO{}

	if err = dtoMapper.Map(companyDTO, company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companyDTO)
}

// GetCompany godoc
// @Summary      Get company details
// @Description  Retrieves the details of a company by its ID
// @Tags         company
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Security     BearerAuth
// @Router       /company/{company_id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	if !h.userCompanyService.UserBelongsToCompany(*userID, uint(companyID)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to get this company's routes"})
		return
	}

	companyDTO := &dtos.CompanyDTO{}

	if err = dtoMapper.Map(companyDTO, company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companyDTO)
}

// UpdateCompany godoc
// @Summary      Update company details
// @Description  Updates the details of an existing company
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        company body CompanyRequest true "Updated company details"
// @Security     BearerAuth
// @Router       /company/{company_id} [put]
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	var companyRequest CompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCompany, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    *userID,
		CompanyID: uint(companyID),
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if userCompany[0].Role != string(RoleAdmin) && userCompany[0].Role != string(RoleManager) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := c.ShouldBindJSON(&companyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	company.Name = companyRequest.Name
	company.Address = companyRequest.Address
	company.Users = nil
	company.Routes = nil
	company.Deliveries = nil
	company.Creator = models.User{}

	if err := h.companyService.Update(context.Background(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	companyDTO := &dtos.CompanyDTO{}
	if err = dtoMapper.Map(companyDTO, company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companyDTO)
}

// DeleteCompany godoc
// @Summary      Delete a company
// @Description  Deletes a company by its ID
// @Tags         company
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Security     BearerAuth
// @Router       /company/{company_id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if company.CreatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.companyService.Delete(context.Background(), uint(companyID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}

// AddUserToCompany godoc
// @Summary      Add a user to a company
// @Description  Adds a user to a company if the request is made by the company creator
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        userID body AddUserToCompanyRequest true "User ID to add"
// @Security     BearerAuth
// @Router       /company/{company_id}/add-user [post]
func (h *CompanyHandler) AddUserToCompany(c *gin.Context) {
	var addUserToCompanyRequest AddUserToCompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if company.CreatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := c.ShouldBindJSON(&addUserToCompanyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userCompany := models.UserCompany{
		UserID:    addUserToCompanyRequest.UserID,
		CompanyID: uint(companyID),
		Role:      string(addUserToCompanyRequest.Role),
	}

	if err := h.userCompanyService.Create(context.Background(), &userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userCompanyDTO := &dtos.UserCompanyDTO{}

	if err = dtoMapper.Map(userCompanyDTO, userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userCompanyDTO)
}

// UpdateUserInCompany godoc
// @Summary      Update a user in a company
// @Description  Updates a user in a company if the request is made by the company creator
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        userID body UpdateUserInCompanyRequest true "User ID to update"
// @Security     BearerAuth
// @Router       /company/{company_id}/update-user [put]
func (h *CompanyHandler) UpdateUserInCompany(c *gin.Context) {
	var updateUserInCompanyRequest UpdateUserInCompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if company.CreatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := c.ShouldBindJSON(&updateUserInCompanyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userCompanies, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    updateUserInCompanyRequest.UserID,
		CompanyID: uint(companyID),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in the company"})
		return
	}

	userCompany := &userCompanies[0]
	if updateUserInCompanyRequest.Role != "" {
		userCompany.Role = updateUserInCompanyRequest.Role
	}

	if err := h.userCompanyService.Update(context.Background(), userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userCompanyDTO := &dtos.UserCompanyDTO{}

	if err = dtoMapper.Map(userCompanyDTO, userCompany); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userCompanyDTO)
}

// RemoveUserFromCompany godoc
// @Summary      Remove a user from a company
// @Description  Removes a user from a company if the request is made by the company creator
// @Tags         company
// @Accept       json
// @Produce      json
// @Param        company_id path int true "Company ID"
// @Param        userID body RemoveUserFromCompanyRequest true "User ID to remove"
// @Security     BearerAuth
// @Router       /company/{company_id}/remove-user [delete]
func (h *CompanyHandler) RemoveUserFromCompany(c *gin.Context) {
	var removeUserFromCompanyRequest RemoveUserFromCompanyRequest

	companyID, err := strconv.Atoi(c.Param("company_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	company, err := h.companyService.GetByID(context.Background(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if company.CreatorID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := c.ShouldBindJSON(&removeUserFromCompanyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userCompanies, err := h.userCompanyService.Where(context.Background(), &models.UserCompany{
		UserID:    removeUserFromCompanyRequest.UserID,
		CompanyID: uint(companyID),
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in the company"})
		return
	}

	userCompany := userCompanies[0]

	if err := h.userCompanyService.Delete(context.Background(), userCompany.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from company successfully"})
}

func getUserIDFromToken(c *gin.Context) (*uint, error) {
	tokenCookie, err := c.Cookie("token")
	if err != nil {
		return nil, errors.New("cookie token not found")
	}

	token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte("mysecret123"), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		stringUserID, ok := claims["sub"].(string)
		if !ok {
			return nil, errors.New("sub not found in token claims")
		}

		userID, err := strconv.Atoi(stringUserID)
		if err != nil {
			return nil, errors.New("problem parsing user ID")
		}

		uintUserID := uint(userID)
		return &uintUserID, nil
	}

	return nil, errors.New("failed to parse token claims")
}
