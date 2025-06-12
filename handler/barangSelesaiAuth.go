package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddPenyelesaianBarangJadi creates a new penyelesaianBarangJadi
func AddPenyelesaianBarangJadi(c *gin.Context) {
	var input model.PenyelesaianBarangJadi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	// Validate required fields
	if input.IDPerintahKerja == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Perintah Kerja is required"})
		return
	}
	if input.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama is required"})
		return
	}

	// Insert the new penyelesaianBarangJadi entry
	query := `
		INSERT INTO "penyelesaianBarangJadi" (id_perintah_kerja, nama, jumlah, tanggal_penyelesaian)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := db.GetDB().QueryRow(query, input.IDPerintahKerja, input.Nama, input.Jumlah, input.TanggalPenyelesaian.ToTime()).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add penyelesaian barang jadi: " + err.Error()})
		return
	}

	// Set the ID in the response
	input.ID = id

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    input,
	})
}

// GetPenyelesaianBarangJadi lists all penyelesaianBarangJadi records
func GetPenyelesaianBarangJadi(c *gin.Context) {
	query := `SELECT id, id_perintah_kerja, nama, jumlah, tanggal_penyelesaian FROM "penyelesaianBarangJadi" ORDER BY id DESC`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch penyelesaian barang jadi: " + err.Error()})
		return
	}
	defer rows.Close()

	var result []model.PenyelesaianBarangJadi
	for rows.Next() {
		var item model.PenyelesaianBarangJadi
		var tanggal sql.NullTime

		if err := rows.Scan(&item.ID, &item.IDPerintahKerja, &item.Nama, &item.Jumlah, &tanggal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse penyelesaian barang jadi: " + err.Error()})
			return
		}

		// Handle null date
		if tanggal.Valid {
			item.TanggalPenyelesaian = model.CustomDate2(tanggal.Time)
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating rows: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result})
}

// GetPenyelesaianBarangJadiByID fetches a single penyelesaianBarangJadi by ID
func GetPenyelesaianBarangJadiByID(c *gin.Context) {
	idStr := c.Param("id")

	// Validate ID parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	query := `SELECT id, id_perintah_kerja, nama, jumlah, tanggal_penyelesaian FROM "penyelesaianBarangJadi" WHERE id=$1`
	row := db.GetDB().QueryRow(query, id)

	var item model.PenyelesaianBarangJadi
	var tanggal sql.NullTime

	err = row.Scan(&item.ID, &item.IDPerintahKerja, &item.Nama, &item.Jumlah, &tanggal)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Penyelesaian barang jadi not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	// Handle null date
	if tanggal.Valid {
		item.TanggalPenyelesaian = model.CustomDate2(tanggal.Time)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    item})
}

// UpdatePenyelesaianBarangJadi updates an existing penyelesaianBarangJadi
func UpdatePenyelesaianBarangJadi(c *gin.Context) {
	idStr := c.Param("id")

	// Validate ID parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	var input model.PenyelesaianBarangJadi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	// Validate required fields
	if input.IDPerintahKerja == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Perintah Kerja is required"})
		return
	}
	if input.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama is required"})
		return
	}

	query := `
		UPDATE "penyelesaianBarangJadi"
		SET id_perintah_kerja = $1, nama = $2, jumlah = $3, tanggal_penyelesaian = $4
		WHERE id = $5
	`
	result, err := db.GetDB().Exec(query, input.IDPerintahKerja, input.Nama, input.Jumlah, input.TanggalPenyelesaian.ToTime(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update penyelesaian barang jadi: " + err.Error()})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check update result: " + err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penyelesaian barang jadi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Penyelesaian Barang Jadi updated successfully"})
}

// DeletePenyelesaianBarangJadi deletes a penyelesaianBarangJadi by ID
func DeletePenyelesaianBarangJadi(c *gin.Context) {
	idStr := c.Param("id")

	// Validate ID parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	query := `DELETE FROM "penyelesaianBarangJadi" WHERE id=$1`
	result, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete penyelesaian barang jadi: " + err.Error()})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check delete result: " + err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Penyelesaian barang jadi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Penyelesaian Barang Jadi deleted successfully"})
}
