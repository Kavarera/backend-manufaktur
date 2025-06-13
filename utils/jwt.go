package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT Claims structure with bitwise roles
type Claims struct {
	Username string `json:"username"`
	Roles    int    `json:"roles"` // Changed from Role string to Roles int
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key") // Change this to environment variable

// GenerateJWTWithRoles generates JWT token with bitwise roles
func GenerateJWTWithRoles(username string, roles int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates JWT token and returns claims
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Backward compatibility function (if you still need it somewhere)
func GenerateJWT(username string, role string) (string, error) {
	// This is for backward compatibility
	// You can map old role strings to bitwise values if needed
	roleMap := map[string]int{
		"SuperAdmin":            63, // All roles
		"BarangManagement":      1,
		"RencanaProduksi":       2,
		"PerintahKerja":         4,
		"HapusPerintahKerja":    8,
		"PengambilanBarangBaku": 16,
		"PengambilanBarangJadi": 32,
	}

	roles, exists := roleMap[role]
	if !exists {
		roles = 1 // Default to BarangManagement
	}

	return GenerateJWTWithRoles(username, roles)
}

// Helper functions for role checking
func HasRole(userRoles int, role int) bool {
	return (userRoles & role) == role
}

// Check if user has any of the specified roles
func HasAnyRole(userRoles int, roles ...int) bool {
	for _, role := range roles {
		if HasRole(userRoles, role) {
			return true
		}
	}
	return false
}

// Check if user has all of the specified roles
func HasAllRoles(userRoles int, roles ...int) bool {
	for _, role := range roles {
		if !HasRole(userRoles, role) {
			return false
		}
	}
	return true
}
