package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	username := c.Param("username")

	query := `
		SELECT ua."id", ua."username", ua."hak_akses"
		FROM "userAccount" ua
		WHERE ua."username" = $1
	`

	row := db.GetDB().QueryRow(query, username)

	var user model.GetUser
	err := row.Scan(&user.UserID, &user.Username, &user.HakAkses)
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
		"message": "User fetched successfully",
		"data":    user,
	})
}

func AllUserList(c *gin.Context) {
	query := `
		SELECT ua."id", ua."username", ua."hak_akses"
		FROM "userAccount" ua
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch users list"})
		return
	}
	defer rows.Close()

	var list []model.GetUser
	var HakAkses int
	for rows.Next() {
		var user model.GetUser
		err := rows.Scan(&user.UserID, &user.Username, &user.HakAkses)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "Error",
				"message": "No Users Available"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Error",
				"message": "Failed to fetch user",
			})
			return
		}
		list = append(list, user)
	}

	// Assuming you have a function `getRoleNames` to fetch role names based on hak_akses values

	roleNames := getRoleNames(HakAkses)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Users fetched successfully",
		"data": gin.H{
			"users": list,
			"roles": roleNames,
		},
	})
}

func UserDelete(c *gin.Context) {
	username := c.Param("username")

	query := `DELETE FROM "userAccount" WHERE "username" = $1`

	res, err := db.GetDB().Exec(query, username)
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
