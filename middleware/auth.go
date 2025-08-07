package middleware

import (
	"manufacture_API/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Role constants using bitwise flags
const (
	RoleBarangManagement      = 1 << 0 // 1 (binary: 000001)
	RoleRencanaProduksi       = 1 << 1 // 2 (binary: 000010)
	RolePerintahKerja         = 1 << 2 // 4 (binary: 000100)
	RoleHapusPerintahKerja    = 1 << 3 // 8 (binary: 001000)
	RolePengambilanBarangBaku = 1 << 4 // 16 (binary: 010000)
	RolePengambilanBarangJadi = 1 << 5 // 32 (binary: 100000)

	// Combined roles
	RoleSuperAdmin = RoleBarangManagement | RoleRencanaProduksi | RolePerintahKerja |
		RoleHapusPerintahKerja | RolePengambilanBarangBaku | RolePengambilanBarangJadi // 63
)

// Permission constants
const (
	// Barang permissions
	PermBarangCreate = "barang:create"
	PermBarangRead   = "barang:read"
	PermBarangUpdate = "barang:update"
	PermBarangDelete = "barang:delete"

	// Gudang permissions
	PermGudangCreate = "gudang:create"
	PermGudangRead   = "gudang:read"
	PermGudangUpdate = "gudang:update"
	PermGudangDelete = "gudang:delete"

	// Mentah permissions
	PermMentahCreate = "mentah:create"
	PermMentahRead   = "mentah:read"
	PermMentahUpdate = "mentah:update"
	PermMentahDelete = "mentah:delete"

	// Rencana permissions
	PermRencanaCreate = "rencana:create"
	PermRencanaRead   = "rencana:read"
	PermRencanaUpdate = "rencana:update"
	PermRencanaDelete = "rencana:delete"
	PermJadwalRead    = "jadwal:read"

	// Perintah permissions
	PermPerintahCreate = "perintah:create"
	PermPerintahRead   = "perintah:read"
	PermPerintahUpdate = "perintah:update"
	PermPerintahDelete = "perintah:delete"

	// Pengambilan permissions
	PermPengambilanCreate = "pengambilan:create"
	PermPengambilanRead   = "pengambilan:read"
	PermPengambilanUpdate = "pengambilan:update"
	PermPengambilanDelete = "pengambilan:delete"

	// Selesai permissions
	PermSelesaiCreate = "selesai:create"
	PermSelesaiRead   = "selesai:read"
	PermSelesaiUpdate = "selesai:update"
	PermSelesaiDelete = "selesai:delete"

	// History Permissions
	PermHistoryRead = "history:read"

	//Super Admin Permissions
	PermUsersCreate = "users:create"
	PermUsersRead   = "users:read"
	PermUsersUpdate = "users:update"
	PermUsersDelete = "users:delete"

	//Formula Permissions
	PermFormulaCreate = "formula:create"
	PermFormulaRead   = "formula:read"
	PermFormulaUpdate = "formula:update"
	PermFormulaDelete = "formula:delete"

	//Barang Satuan Permissionns
	PermSatuanCreate = "satuan:create"
	PermSatuanRead   = "satuan:read"
	PermSatuanUpdate = "satuan:update"
	PermSatuanDelete = "satuan:delete"
)

// Helper functions for role management
func HasRole(userRoles int, role int) bool {
	return (userRoles & role) == role
}

func AddRole(userRoles int, role int) int {
	return userRoles | role
}

func RemoveRole(userRoles int, role int) int {
	return userRoles &^ role
}

func GetRoleNames(roles int) []string {
	var roleNames []string

	if HasRole(roles, RoleBarangManagement) {
		roleNames = append(roleNames, "BarangManagement")
	}
	if HasRole(roles, RoleRencanaProduksi) {
		roleNames = append(roleNames, "RencanaProduksi")
	}
	if HasRole(roles, RolePerintahKerja) {
		roleNames = append(roleNames, "PerintahKerja")
	}
	if HasRole(roles, RoleHapusPerintahKerja) {
		roleNames = append(roleNames, "HapusPerintahKerja")
	}
	if HasRole(roles, RolePengambilanBarangBaku) {
		roleNames = append(roleNames, "PengambilanBarangBaku")
	}
	if HasRole(roles, RolePengambilanBarangJadi) {
		roleNames = append(roleNames, "PengambilanBarangJadi")
	}

	return roleNames
}

// checkUserPermission checks if user roles have specific permission
func checkUserPermission(userRoles int, requiredPermission string) bool {
	// Define role-permission mapping
	rolePermissions := map[int][]string{
		RoleBarangManagement: {
			PermBarangCreate, PermBarangRead, PermBarangUpdate, PermBarangDelete,
			PermGudangCreate, PermGudangRead, PermGudangUpdate, PermGudangDelete,
			PermMentahCreate, PermMentahRead, PermMentahUpdate, PermMentahDelete,
			PermFormulaCreate, PermFormulaRead, PermFormulaUpdate, PermFormulaDelete,
			PermSatuanCreate, PermSatuanRead, PermSatuanUpdate, PermSatuanDelete, PermHistoryRead,
		},
		RoleRencanaProduksi: {
			PermRencanaCreate, PermRencanaRead, PermRencanaUpdate, PermRencanaDelete, PermJadwalRead,
			PermFormulaCreate, PermFormulaRead, PermFormulaUpdate, PermFormulaDelete, PermHistoryRead,
		},
		RolePerintahKerja: {
			PermPerintahCreate, PermPerintahRead, PermPerintahUpdate, PermHistoryRead, PermRencanaRead,
		},
		RoleHapusPerintahKerja: {
			PermPerintahDelete, PermHistoryRead, PermRencanaRead,
		},
		RolePengambilanBarangBaku: {
			PermPengambilanCreate, PermPengambilanRead, PermPengambilanUpdate, PermPengambilanDelete, PermHistoryRead,
		},
		RolePengambilanBarangJadi: {
			PermSelesaiCreate, PermSelesaiRead, PermSelesaiUpdate, PermSelesaiDelete, PermHistoryRead, PermPerintahRead,
		},
	}

	// Check each role the user has
	for role, permissions := range rolePermissions {
		if HasRole(userRoles, role) {
			for _, permission := range permissions {
				if permission == requiredPermission {
					return true
				}
			}
		}
	}

	// Super Admin has all permissions
	if userRoles == RoleSuperAdmin {
		return true
	}

	return false
}

// BitwisePermissionMiddleware checks if user has specific permission using bitwise roles
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "Error",
				"message": "Missing authorization token",
			})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Validate token and extract claims
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "Error",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if user has the required permission
		// Note: You'll need to modify your JWT claims to include user roles as int
		hasPermission := checkUserPermission(claims.Roles, requiredPermission)

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"status":              "Error",
				"message":             "Insufficient permissions for this action",
				"required_permission": requiredPermission,
			})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("username", claims.Username)
		c.Set("user_roles", claims.Roles)
		c.Set("user_role_names", GetRoleNames(claims.Roles))

		c.Next()
	}
}

// RoleMiddleware checks if user has specific role(s)
func RoleMiddleware(requiredRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "Error",
				"message": "Missing authorization token",
			})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "Error",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range requiredRoles {
			if HasRole(claims.Roles, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "Error",
				"message": "Insufficient role permissions",
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("user_roles", claims.Roles)
		c.Set("user_role_names", GetRoleNames(claims.Roles))

		c.Next()
	}
}
