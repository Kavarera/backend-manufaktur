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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user model.User
	query := `
		SELECT u."id", u."username", u."password", h."hak_akses", u."hak_akses"
		FROM "userAccount" u
		JOIN "hakAkses" h ON u."hak_akses" = h."id"
		WHERE u."username" = $1
	`

	err := db.GetDB().QueryRow(query, credentials.Username).Scan(
		&user.UserID, &user.Username, &user.Password, &user.HakAkses, &user.IdHakAkses,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Match hashed password
	hashedInput := hashPassword(credentials.Password)
	if user.Password != hashedInput {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password Salah"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.Username, user.HakAkses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"userId":   user.UserID,
			"username": user.Username,
			"role":     user.HakAkses,
			"roleId":   user.IdHakAkses,
		},
	})
}

// Utility function to hash password
func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
