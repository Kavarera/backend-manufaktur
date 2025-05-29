package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	var user struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		IdHakAkses string `json:"hak_akses"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid request",
		})
		return
	}

	// Generate UUID for user ID
	userID := uuid.New().String()

	roleID, err := strconv.Atoi(user.IdHakAkses)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid role ID format",
		})
		return
	}

	var roleName string
	query := `SELECT "hak_akses" FROM "hakAkses" WHERE "id" = $1`
	err = db.GetDB().QueryRow(query, roleID).Scan(&roleName)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid role ID: Role does not exist",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to validate role",
		})
		return
	}

	query = `SELECT "username" FROM "userAccount" WHERE "username" = $1`
	var existingUsername string
	err = db.GetDB().QueryRow(query, user.Username).Scan(&existingUsername)
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

	query = `
		INSERT INTO "userAccount" ("id", "username", "password", "hak_akses")
		VALUES ($1, $2, $3, $4)
	`
	_, err = db.GetDB().Exec(query, userID, user.Username, hashedPassword, roleID)
	if err != nil {
		fmt.Printf("Database insert error: %v\n", err)
		fmt.Printf("UserID: %s, Username: %s, RoleID: %d\n", userID, user.Username, roleID)

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Failed to register user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data": gin.H{
			"userId":   userID,
			"username": user.Username,
			"roleId":   roleID,
		},
	})
}
