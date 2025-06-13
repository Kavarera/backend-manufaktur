package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"manufacture_API/db"
	"manufacture_API/model"
	"manufacture_API/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid request",
		})
		return
	}

	var user model.User
	// Updated query to get roles as integer
	query := `
		SELECT "id", "username", "password", "hak_akses"
		FROM "userAccount"
		WHERE "username" = $1
	`

	err := db.GetDB().QueryRow(query, credentials.Username).Scan(
		&user.UserID, &user.Username, &user.Password, &user.HakAkses,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Error",
			"message": "Invalid credentials",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Database error",
		})
		return
	}

	hashedInput := hashPassword(credentials.Password)
	if user.Password != hashedInput {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "Error",
			"message": "Invalid credentials",
		})
		return
	}

	// Generate JWT token with roles as integer
	token, err := utils.GenerateJWTWithRoles(user.Username, user.HakAkses) // You'll need to update this function
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Could not generate token",
		})
		return
	}

	// Get role names for response
	roleNames := getRoleNames(user.HakAkses)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Login successful",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"userId":    user.UserID,
				"username":  user.Username,
				"roles":     roleNames,
				"roleValue": user.HakAkses,
			},
		},
	})
}

// Helper function to get role names from bitwise value
func getRoleNames(roles int) []string {
	var roleNames []string

	if (roles & RoleBarangManagement) != 0 {
		roleNames = append(roleNames, "BarangManagement")
	}
	if (roles & RoleRencanaProduksi) != 0 {
		roleNames = append(roleNames, "RencanaProduksi")
	}
	if (roles & RolePerintahKerja) != 0 {
		roleNames = append(roleNames, "PerintahKerja")
	}
	if (roles & RoleHapusPerintahKerja) != 0 {
		roleNames = append(roleNames, "HapusPerintahKerja")
	}
	if (roles & RolePengambilanBarangBaku) != 0 {
		roleNames = append(roleNames, "PengambilanBarangBaku")
	}
	if (roles & RolePengambilanBarangJadi) != 0 {
		roleNames = append(roleNames, "PengambilanBarangJadi")
	}

	// Check if it's SuperAdmin (has all roles)
	if roles == RoleSuperAdmin {
		return []string{"SuperAdmin"}
	}

	return roleNames
}

// UpdateUserRoles - New endpoint to update user roles dynamically
func UpdateUserRoles(c *gin.Context) {
	userName := c.Param("username")

	var updateData struct {
		Roles []string `json:"roles"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid request",
		})
		return
	}

	// Parse new roles
	newRoles := parseRolesFromRequest(updateData.Roles)
	if newRoles == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid roles provided",
		})
		return
	}

	// Update user roles in database
	query := `UPDATE "userAccount" SET "hak_akses" = $1 WHERE "username" = $2`
	result, err := db.GetDB().Exec(query, newRoles, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to update user roles",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User roles updated successfully",
		"data": gin.H{
			"userId":    userName,
			"roles":     updateData.Roles,
			"roleValue": newRoles,
		},
	})
}

// GetUserRoles - Get current user roles
func GetUserRoles(c *gin.Context) {
	userNAME := c.Param("username")

	var userRoles int
	query := `SELECT "hak_akses" FROM "userAccount" WHERE "username" = $1`
	err := db.GetDB().QueryRow(query, userNAME).Scan(&userRoles)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "User not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Database error",
		})
		return
	}

	roleNames := getRoleNames(userRoles)

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data": gin.H{
			"userId":    userNAME,
			"roles":     roleNames,
			"roleValue": userRoles,
		},
	})
}

// Utility function to hash password (keep your existing implementation)
func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
