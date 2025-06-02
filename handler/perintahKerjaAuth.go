package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListPerintahKerja returns all records
func ListPerintahKerja(c *gin.Context) {
	query := `
		SELECT "id", "tanggal_rilis", "tanggal_progres", "tanggal_selesai", 
		       "status", "hasil", "customer", "keterangan", "document_url", "document_nama"
		FROM "perintahKerja" 
		ORDER BY "id"
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to fetch data: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var list []model.PerintahKerja
	for rows.Next() {
		var data model.PerintahKerja
		err := rows.Scan(
			&data.ID, &data.TanggalRilis, &data.TanggalProgres, &data.TanggalSelesai,
			&data.Status, &data.Hasil, &data.Customer, &data.Keterangan,
			&data.DocumentURL, &data.DocumentNama,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Error",
				"message": "Failed to parse data: " + err.Error(),
			})
			return
		}
		list = append(list, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    list,
	})
}

// AddPerintahKerja adds a new record
func AddPerintahKerja(c *gin.Context) {
	var input model.PerintahKerja
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	// Generate ID if not provided
	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	query := `
		INSERT INTO "perintahKerja" 
		("id", "tanggal_rilis", "tanggal_progres", "tanggal_selesai", "status", "hasil", "customer", "keterangan", "document_url", "document_nama")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING "id"
	`

	var returnedID string
	err := db.GetDB().QueryRow(query,
		input.ID, input.TanggalRilis, input.TanggalProgres, input.TanggalSelesai,
		input.Status, input.Hasil, input.Customer, input.Keterangan,
		input.DocumentURL, input.DocumentNama,
	).Scan(&returnedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to insert data: " + err.Error(),
		})
		return
	}

	input.ID = returnedID
	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Berhasil menambah perintah kerja",
		"data":    input,
	})
}

// UpdatePerintahKerja updates by ID
func UpdatePerintahKerja(c *gin.Context) {
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

	if pk.TanggalRilis != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_rilis"=$%d`, argPos))
		values = append(values, pk.TanggalRilis)
		argPos++
	}
	if pk.TanggalProgres != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_progres"=$%d`, argPos))
		values = append(values, pk.TanggalProgres)
		argPos++
	}
	if pk.TanggalSelesai != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_selesai"=$%d`, argPos))
		values = append(values, pk.TanggalSelesai)
		argPos++
	}
	if pk.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf(`"status"=$%d`, argPos))
		values = append(values, pk.Status)
		argPos++
	}
	if pk.Hasil != 0 {
		setClauses = append(setClauses, fmt.Sprintf(`"hasil"=$%d`, argPos))
		values = append(values, pk.Hasil)
		argPos++
	}
	if pk.Customer != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"customer"=$%d`, argPos))
		values = append(values, pk.Customer)
		argPos++
	}
	if pk.Keterangan != "" {
		setClauses = append(setClauses, fmt.Sprintf(`"keterangan"=$%d`, argPos))
		values = append(values, pk.Keterangan)
		argPos++
	}
	if pk.DocumentURL != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"document_url"=$%d`, argPos))
		values = append(values, pk.DocumentURL)
		argPos++
	}
	if pk.DocumentNama != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"document_nama"=$%d`, argPos))
		values = append(values, pk.DocumentNama)
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
			"message": "Failed to update perintah kerja: " + err.Error(),
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
		"message": "Perintah kerja updated successfully",
	})
}

// DeletePerintahKerja deletes by ID
func DeletePerintahKerja(c *gin.Context) {
	id := c.Param("id")

	// First, get document info if exists
	var documentURL sql.NullString
	err := db.GetDB().QueryRow(`SELECT "document_url" FROM "perintahKerja" WHERE "id"=$1`, id).Scan(&documentURL)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to check document info",
		})
		return
	}

	// Delete from database
	res, err := db.GetDB().Exec(`DELETE FROM "perintahKerja" WHERE "id"=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to delete data: " + err.Error(),
		})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Data not found",
		})
		return
	}

	// Delete associated file if exists
	if documentURL.Valid && documentURL.String != "" {
		// Extract filename from URL
		parts := strings.Split(documentURL.String, "/")
		if len(parts) > 0 {
			filename := parts[len(parts)-1]
			filePath := filepath.Join("./uploads/documents", filename)
			os.Remove(filePath) // Ignore error if file doesn't exist
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Data berhasil dihapus",
	})
}
