package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT ua."id", ua."username", ua."password", ha."hak_akses", ua."hak_akses"
		FROM "userAccount" ua
		JOIN "hakAkses" ha ON ua."hak_akses" = ha."id"
		WHERE ua."id" = $1
	`

	row := db.GetDB().QueryRow(query, id)

	var user model.User
	err := row.Scan(&user.UserID, &user.Username, &user.Password, &user.HakAkses, &user.IdHakAkses)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM "userAccount" WHERE id = $1`

	_, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
