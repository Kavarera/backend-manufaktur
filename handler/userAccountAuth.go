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
		SELECT ua."id", ua."username", ha."hak_akses", ua."hak_akses"
		FROM "userAccount" ua
		JOIN "hakAkses" ha ON ua."hak_akses" = ha."id"
		WHERE ua."id" = $1
	`

	row := db.GetDB().QueryRow(query, id)

	var user model.GetUser
	err := row.Scan(&user.UserID, &user.Username, &user.HakAkses, &user.IdHakAkses)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "User not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to fetch user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    user,
	})
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM "userAccount" WHERE id = $1`

	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to delete user",
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "User deleted successfully",
	})
}
