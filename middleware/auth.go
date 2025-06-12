package middleware

import (
	"manufacture_API/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware checks if user has specific permission
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Validate token and extract claims
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		// Assuming you have permissions stored in claims or you fetch them based on role
		hasPermission := checkUserPermission(claims.Role, requiredPermission)

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// checkUserPermission checks if a role has specific permission
func checkUserPermission(userRole, requiredPermission string) bool {
	// Define role-permission mapping
	rolePermissions := map[string][]string{
		"SuperAdmin": {
			"users:create", "users:read", "users:update", "users:delete",
			"barang:create", "barang:read", "barang:update", "barang:delete",
			"gudang:create", "gudang:read", "gudang:update", "gudang:delete",
			"mentah:create", "mentah:read", "mentah:update", "mentah:delete",
			"rencana:create", "rencana:read", "rencana:update", "rencana:delete", "jadwal:read",
			"perintah:create", "perintah:read", "perintah:update", "perintah:delete",
			"pengambilan:create", "pengambilan:read", "pengambilan:update", "pengambilan:delete",
			"selesai:create", "pselesai:read", "selesai:update", "selesai:delete", "history:read",
		},
		"BarangManagement": {
			"barang:create", "barang:read", "barang:update", "barang:delete",
			"gudang:create", "gudang:read", "gudang:update", "gudang:delete",
			"mentah:create", "mentah:read", "mentah:update", "mentah:delete",
		},
		"RencanaProduksi": {
			"rencana:create", "rencana:read", "rencana:update", "rencana:delete", "jadwal:read",
		},
		"PerintahKerja": {
			"perintah:create", "perintah:read", "perintah:update", "history:read",
		},
		"HapusPerintahKerja": {
			"perintah:delete",
		},
		"PengambilanBarangBaku": {
			"pengambilan:create", "pengambilan:read", "pengambilan:update", "pengambilan:delete",
		},
		"PengambilanBarangJadi": {
			"selesai:create", "pselesai:read", "selesai:update", "selesai:delete",
		},
	}

	permissions, exists := rolePermissions[userRole]
	if !exists {
		return false
	}

	for _, permission := range permissions {
		if permission == requiredPermission {
			return true
		}
	}
	return false
}
