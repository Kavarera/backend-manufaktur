package handler

import (
	"API/db"
	"API/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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

	// Database table selection
	var user mode.User
	query := `
        SELECT u."id", u."username", u."password", u."hak_akses",h."hak_akses"
        FROM "userAccount" u
        JOIN "hakAkses" h ON u."hak_akses" = h."id"
        WHERE u."username" = $1;
    `
	err := db.DB.QueryRow(query, credentials.Username).Scan(
		&user.userID, &user.Username, &user.Password, &user.hakAkses,
	)
	if err == sql.ErrNoRows || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	hashedInputPassword := hashPassword(credentials.Password)
	if hashedInputPassword != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.Username, user.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"username": user.Username,
			"rolename": user.RoleName,
			"fullname": user.Fullname,
		},
	})
}

// Utility function to hash password
func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
