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
		Username   string `json:"username"`  // Required field for unique username
		Password   string `json:"password"`  // Required field for user's password
		IdHakAkses string `json:"hak_akses"` // Re	quired field to assign a role
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate role_id
	roleID, err := strconv.Atoi(user.IdHakAkses)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID format"})
		return
	}

	// Check if role_id exists in the role_table
	var roleName string
	query := `SELECT "hak_akses" FROM "hakAkses" WHERE "id" = $1`
	err = db.DB.QueryRow(query, roleID).Scan(&roleName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID: Role does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate role"})
		}
		return
	}

	// Check if username is unique
	query = `SELECT "username" FROM "userAccount" WHERE "username" = $1`
	var existingUsername string
	err = db.DB.QueryRow(query, user.Username).Scan(&existingUsername)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username availability"})
		return
	}

	// Hash the password
	hashedPassword := hashPassword(user.Password)

	// Insert the new user into the database
	query = `
        INSERT INTO "userAccount" ("id", "username", "password", "hak_akses")
        VALUES ($1, $2, $3, $4)
    `
	_, err = db.DB.Exec(query, user.Id, user.Username, hashedPassword, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}
