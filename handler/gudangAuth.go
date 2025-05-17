package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListGudang returns all gudang records
func ListGudang(c *gin.Context) {
	rows, err := db.GetDB().Query(`SELECT "id", "nama" FROM "gudang" ORDER BY "id"`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gudang list"})
		return
	}
	defer rows.Close()

	var list []model.Gudang
	for rows.Next() {
		var g model.Gudang
		if err := rows.Scan(&g.ID, &g.Nama); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse gudang data"})
			return
		}
		list = append(list, g)
	}

	c.JSON(http.StatusOK, list)
}

// GetGudangByID returns a single gudang by ID
func GetGudangByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	row := db.GetDB().QueryRow(`SELECT "id", "nama" FROM "gudang" WHERE "id" = $1`, id)
	var g model.Gudang
	err = row.Scan(&g.ID, &g.Nama)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gudang not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gudang"})
		}
		return
	}

	c.JSON(http.StatusOK, g)
}

// AddGudang inserts a new gudang record
func AddGudang(c *gin.Context) {
	var g model.Gudang

	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `INSERT INTO "gudang" ("nama") VALUES ($1) RETURNING "id"`
	err := db.GetDB().QueryRow(query, g.Nama).Scan(&g.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gudang"})
		return
	}

	c.JSON(http.StatusCreated, g)
}

// UpdateGudang updates the nama of a gudang by ID
func UpdateGudang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var g model.Gudang
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `UPDATE "gudang" SET "nama"=$1 WHERE "id"=$2`
	res, err := db.GetDB().Exec(query, g.Nama, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update gudang"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gudang not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gudang updated successfully"})
}

// DeleteGudang deletes a gudang by ID
func DeleteGudang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	query := `DELETE FROM "gudang" WHERE "id"=$1`
	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete gudang"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gudang not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gudang deleted successfully"})
}
