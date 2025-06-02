package handler

import (
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateProsesPengerjaan(c *gin.Context) {
	id := c.Param("id")

	var pk model.PerintahKerja
	if err := c.ShouldBindJSON(&pk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	setClauses := []string{}
	values := []interface{}{}
	argPos := 1

	if pk.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf(`"status"=$%d`, argPos))
		values = append(values, pk.Status)
		argPos++
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "No fields to update",
		})
		return
	}

	values = append(values, id)
	query := fmt.Sprintf(`UPDATE "perintahKerja" SET %s WHERE "id"=$%d`, strings.Join(setClauses, ", "), argPos)

	res, err := db.GetDB().Exec(query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to update status perintah kerja: " + err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Perintah kerja not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Status perintah kerja updated successfully",
	})
}
