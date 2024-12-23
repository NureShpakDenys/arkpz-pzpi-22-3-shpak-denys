package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
)

// AdminHandler is a handler for admin endpoints
type AdminHandler struct {
	dbPassword    string               // Password for the database
	encryptionKey string               // Encryption key for sensitive data
	userService   services.UserService // Service for user operations
}

// Config holds configuration data loaded from YAML
type Config struct {
	EncryptionKey string `yaml:"encryption_key"`
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(dbPassword string, userService services.UserService, encryptionKey string) *AdminHandler {
	return &AdminHandler{
		dbPassword:    dbPassword,
		userService:   userService,
		encryptionKey: encryptionKey,
	}
}

// encryptData encrypts the given data using the provided key
func encryptData(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptData decrypts the given data using the provided key
func decryptData(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// findLatestFile finds the latest file in the given directory
func findLatestFile(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	var latestFile string
	var latestTime time.Time

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dir, file.Name())
			if file.ModTime().After(latestTime) {
				latestTime = file.ModTime()
				latestFile = filePath
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("no files found in directory")
	}

	return latestFile, nil
}

// BackupDatabase godoc
// @Summary Backup the database
// @Description Creates a backup of the database and encrypts it
// @Produce json
// @Security     BearerAuth
// @Router /admin/backup [post]
func (h *AdminHandler) BackupDatabase(c *gin.Context) {
	backupDir := filepath.Join("../../migrations")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%d.sql", 12345)) // Для примера PID = 12345

	// Преобразование в абсолютный путь
	absPath, err := filepath.Abs(backupPath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	fmt.Println("backupPath: ", absPath)

	file, err := os.Create(absPath)
	if err != nil {
		fmt.Println("Error when creating file:", err)
		return
	}
	file.Close()

	cmd := exec.Command(
		"pg_dump",
		"--dbname=Wayra",
		"--schema=public",
		fmt.Sprintf("--file=%s", backupPath),
		"--data-only",
		"--format=c",
		"--username=postgres",
		"--host=localhost",
		"--port=5432",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+h.dbPassword)

	if err := cmd.Run(); err != nil {
		fmt.Println("cmd run: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create backup: %s", err.Error())})
		return
	}

	backupData, err := ioutil.ReadFile(backupPath)
	if err != nil {
		fmt.Println("ioutil read file: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read backup file: %s", err.Error())})
		return
	}

	encryptedData, err := encryptData(string(backupData), h.encryptionKey)
	if err != nil {
		fmt.Println("encrypt: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to encrypt backup: %s", err.Error())})
		return
	}

	err = ioutil.WriteFile(backupPath, []byte(encryptedData), 0644)
	if err != nil {
		fmt.Println("ioutil.WriteFile: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save encrypted backup: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup created and encrypted successfully", "path": backupPath})
}




// RestoreDatabase godoc
// @Summary Restore the database
// @Description Restores the database from an encrypted backup
// @Produce json
// @Security     BearerAuth
// @Router /admin/restore [post]
func (h *AdminHandler) RestoreDatabase(c *gin.Context) {
	backupDir := filepath.Join("../../migrations")

	absPath, err := filepath.Abs(backupDir)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	latestFile, err := findLatestFile(absPath)
	if err != nil {
		fmt.Println("Error finding latest file:", err)
		return
	}
	fmt.Println("Latest backup file: ", latestFile)

	backupData, err := ioutil.ReadFile(latestFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read backup file: %s", err.Error())})
		return
	}

	decryptedData, err := decryptData(string(backupData), h.encryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decrypt backup: %s", err.Error())})
		return
	}

	if err := ioutil.WriteFile(latestFile, []byte(decryptedData), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save decrypted backup: %s", err.Error())})
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
		latestFile,
	)

	restoreCmd.Stdout = os.Stdout
	restoreCmd.Stderr = os.Stderr
	restoreCmd.Env = append(os.Environ(), "PGPASSWORD="+h.dbPassword)

	if err := restoreCmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to restore database: %s", err.Error())})
		return
	}

	encryptedData, err := encryptData(string(backupData), h.encryptionKey)
	if err != nil {
		fmt.Println("encrypt: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to encrypt backup: %s", err.Error())})
		return
	}

	err = ioutil.WriteFile(latestFile, []byte(encryptedData), 0644)
	if err != nil {
		fmt.Println("ioutil.WriteFile: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save encrypted backup: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database restored successfully"})
}
