package handler

import (
	"database/sql"
	"manufacture_API/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Role constants (same as in middleware)
const (
	RoleBarangManagement      = 1 << 0 // 1
	RoleRencanaProduksi       = 1 << 1 // 2
	RolePerintahKerja         = 1 << 2 // 4
	RoleHapusPerintahKerja    = 1 << 3 // 8
	RolePengambilanBarangBaku = 1 << 4 // 16
	RolePengambilanBarangJadi = 1 << 5 // 32
	RoleSuperAdmin            = 63     // All roles combined
)

// Helper function to parse roles from frontend checkboxes
func parseRolesFromRequest(roles []string) int {
	roleMap := map[string]int{
		"BarangManagement":      RoleBarangManagement,
		"RencanaProduksi":       RoleRencanaProduksi,
		"PerintahKerja":         RolePerintahKerja,
		"HapusPerintahKerja":    RoleHapusPerintahKerja,
		"PengambilanBarangBaku": RolePengambilanBarangBaku,
		"PengambilanBarangJadi": RolePengambilanBarangJadi,
		"SuperAdmin":            RoleSuperAdmin,
	}

	totalRoles := 0
	for _, roleName := range roles {
		if roleValue, exists := roleMap[roleName]; exists {
			totalRoles |= roleValue
		}
	}

	return totalRoles
}

func Register(c *gin.Context) {
	var user struct {
		Username string   `json:"username"`
		Password string   `json:"password"`
		Roles    []string `json:"roles"` // Changed to array of role names
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid request",
		})
		return
	}

	// Validate input
	if user.Username == "" || user.Password == "" || len(user.Roles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Username, password, and roles are required",
		})
		return
	}

	// Generate UUID for user ID
	userID := uuid.New().String()

	// Parse roles to bitwise integer
	userRoles := parseRolesFromRequest(user.Roles)
	if userRoles == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid roles provided",
		})
		return
	}

	// Check if username already exists
	query := `SELECT "username" FROM "userAccount" WHERE "username" = $1`
	var existingUsername string
	err := db.GetDB().QueryRow(query, user.Username).Scan(&existingUsername)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Username already exists",
		})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to check username availability",
		})
		return
	}

	hashedPassword := hashPassword(user.Password)

	// Updated query to store roles as integer
	query = `
		INSERT INTO "userAccount" ("id", "username", "password", "hak_akses")
		VALUES ($1, $2, $3, $4)
	`
	_, err = db.GetDB().Exec(query, userID, user.Username, hashedPassword, userRoles)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Failed to register user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User registered successfully",
		"data": gin.H{
			"userId":    userID,
			"username":  user.Username,
			"roles":     user.Roles,
			"roleValue": userRoles,
		},
	})
}
