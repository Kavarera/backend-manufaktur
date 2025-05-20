package handler

import (
	"database/sql"
	"manufacture_API/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user struct {
		Id         string `json:"id"`
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
	_, err = db.GetDB().Exec(query, user.Id, user.Username, hashedPassword, roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data": gin.H{
			"userId":   user.Id,
			"username": user.Username,
			"roleId":   roleID,
		},
	})
}
