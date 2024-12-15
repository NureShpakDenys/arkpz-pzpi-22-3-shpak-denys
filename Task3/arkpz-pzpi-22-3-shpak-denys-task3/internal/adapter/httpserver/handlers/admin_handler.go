package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
)

// AdminHandler is a handler for admin endpoints
type AdminHandler struct {
	dbPassword  string               // Password for the database
	userService services.UserService // Service for user operations
}

// NewAdminHandler creates a new AdminHandler
// dbPassword: Password for the database
// userService: Service for user operations
// Returns: A new AdminHandler
func NewAdminHandler(dbPassword string, userService services.UserService) *AdminHandler {
	return &AdminHandler{
		dbPassword:  dbPassword,
		userService: userService,
	}
}

// BackupDatabaseRequest is the request for the BackupDatabase endpoint
type BackupDatabaseRequest struct {
	// Backup path is the path where the backup will be stored
	BackupPath string `json:"backup_path"`
}

// BackupDatabase godoc
// @Summary Backup the database
// @Description Backup the database
// @Tags admin
// @Accept json
// @Produce json
// @Param BackupDatabaseRequest body BackupDatabaseRequest true "Backup path"
// @Security     BearerAuth
// @Success 200 {string} string "Backup created"
// @Router /admin/backup [post]
func (h *AdminHandler) BackupDatabase(c *gin.Context) {
	var req BackupDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetByID(context.Background(), *userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.Role.Name != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cmd := exec.Command(
		"pg_dump",
		"--dbname=Wayra",
		"--schema=public",
		fmt.Sprintf("--file=%s", req.BackupPath),
		"--data-only",
		"--format=c",
		"--username=postgres",
		"--host=localhost",
		"--port=5432",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+h.dbPassword)

	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Backup created")
}

// RestoreDatabase godoc
// @Summary Restore the database
// @Description Restore the database
// @Tags admin
// @Accept json
// @Produce json
// @Param BackupDatabaseRequest body BackupDatabaseRequest true "Backup path"
// @Security     BearerAuth
// @Success 200 {string} string "Database restored"
// @Router /admin/restore [post]
func (h *AdminHandler) RestoreDatabase(c *gin.Context) {
	var req BackupDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error getting request": err.Error()})
		return
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetByID(context.Background(), *userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getting user by id": err.Error()})
		return
	}

	if user.Role.Name != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	truncateCmd := exec.Command(
		"psql",
		"--dbname=Wayra",
		"--username=postgres",
		"--host=localhost",
		"--port=5432",
		"-c",
		"TRUNCATE TABLE roles, users, companies, routes, deliveries, product_categories, products, waypoints, sensor_data, user_companies RESTART IDENTITY CASCADE;",
	)
	truncateCmd.Env = append(os.Environ(), "PGPASSWORD="+h.dbPassword)

	if err := truncateCmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while truncating": err.Error()})
		return
	}

	restoreCmd := exec.Command(
		"pg_restore",
		"--no-owner",
		"--role=postgres",
		"--dbname=Wayra",
		"--format=c",
		"-v",
		"--clean",
		"--if-exists",
		"--host=localhost",
		"--port=5432",
		"--username=postgres",
		req.BackupPath,
	)
	restoreCmd.Stdout = os.Stdout
	restoreCmd.Stderr = os.Stderr
	restoreCmd.Env = append(os.Environ(), "PGPASSWORD="+h.dbPassword)
	if err := restoreCmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error while restoring": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Database restored")
}
