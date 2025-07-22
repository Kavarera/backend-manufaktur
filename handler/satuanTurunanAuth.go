package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// List all satuan turunan
func ListSatuanTurunan(c *gin.Context) {
	query := `SELECT id, nama, satuan FROM "satuanTurunan" ORDER BY id`
	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch satuan turunan"})
		return
	}
	defer rows.Close()

	var result []model.SatuanTurunan
	for rows.Next() {
		var item model.SatuanTurunan
		if err := rows.Scan(&item.IDTurunan, &item.NamaSatuanTurunan, &item.IDSatuan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse satuan turunan"})
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

// Get satuan turunan by ID
func GetSatuanTurunanByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var item model.SatuanTurunan
	query := `SELECT id, nama, satuan FROM "satuanTurunan" WHERE id = $1`
	err = db.GetDB().QueryRow(query, id).Scan(&item.IDTurunan, &item.NamaSatuanTurunan, &item.IDSatuan)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Satuan Turunan not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch satuan turunan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    item,
	})
}

// Add multiple satuan turunan
func AddSatuanTurunan(c *gin.Context) {
	var inputs []model.SatuanTurunan
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	stmt := `INSERT INTO "satuanTurunan" ("nama", "satuan") VALUES ($1, $2) RETURNING "id"`

	for i := range inputs {
		err := tx.QueryRow(stmt, inputs[i].NamaSatuanTurunan, inputs[i].IDSatuan).Scan(&inputs[i].IDTurunan)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert satuan turunan"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Satuan Turunan added successfully",
		"data":    inputs,
	})
}

func UpdateSatuanTurunan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	allowedFields := map[string]string{
		"namaTurunan": "nama",
		"idSatuan":    "satuan",
	}

	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for field, value := range input {
		column, ok := allowedFields[field]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field '%s' is not allowed for update", field)})
			return
		}
		setClauses = append(setClauses, fmt.Sprintf(`%s = $%d`, column, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`UPDATE "satuanTurunan" SET %s WHERE id = $%d`, strings.Join(setClauses, ", "), argIndex)
	args = append(args, id)

	_, err = db.GetDB().Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update satuan turunan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Satuan Turunan updated successfully",
	})
}

// Delete satuan turunan by ID
func DeleteSatuanTurunan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `DELETE FROM "satuanTurunan" WHERE id = $1`
	_, err = db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete satuan turunan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Satuan Turunan deleted successfully",
	})
}
