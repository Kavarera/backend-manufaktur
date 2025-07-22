package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListBarangSatuan handles GET /barangSatuan
func ListBarangSatuan(c *gin.Context) {
	query := `SELECT id, nama FROM "barangSatuan" ORDER BY id`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch satuan"})
		return
	}
	defer rows.Close()

	var result []model.Satuan
	for rows.Next() {
		var item model.Satuan
		if err := rows.Scan(&item.IDSatuan, &item.NamaSatuan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse satuan"})
			return
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}

// GetBarangSatuanByID handles GET /barangSatuan/:id
func GetBarangSatuanByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var satuan model.Satuan
	query := `SELECT id, nama FROM "barangSatuan" WHERE id = $1`

	err = db.GetDB().QueryRow(query, id).Scan(&satuan.IDSatuan, &satuan.NamaSatuan)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang Satuan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch satuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    satuan,
	})
}

// AddBarangSatuan handles POST /barangSatuan
func AddBarangSatuan(c *gin.Context) {
	var inputs []model.Satuan
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if len(inputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Any Input Found"})
	}

	dbConn := db.GetDB()
	tx, err := dbConn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Connection Interupted"})
		return
	}

	query := `INSERT INTO "barangSatuan" ("nama") VALUES ($1) RETURNING "id"`
	for i := range inputs {
		err := db.GetDB().QueryRow(query, inputs[i].NamaSatuan).Scan(&inputs[i].IDSatuan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert satuan"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Barang Satuan added successfully",
		"data":    inputs,
	})
}

// UpdateBarangSatuan handles PUT /barangSatuan/:id
func UpdateBarangSatuan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input model.Satuan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	query := `UPDATE "barangSatuan" SET nama = $1 WHERE id = $2`
	_, err = db.GetDB().Exec(query, input.NamaSatuan, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update satuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Barang Satuan updated successfully",
	})
}

// DeleteBarangSatuan handles DELETE /barangSatuan/:id
func DeleteBarangSatuan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `DELETE FROM "barangSatuan" WHERE id = $1`
	_, err = db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete satuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Barang Satuan deleted successfully",
	})
}
