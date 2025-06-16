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

// ListFormulaProduksi lists all formulaProduksi records
func ListFormulaProduksi(c *gin.Context) {
	query := `
		SELECT id, barang_jadi, kuantitas, satuan, bahan_baku
		FROM "formulaProduksi"
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch formula produksi"})
		return
	}
	defer rows.Close()

	var result []model.FormulaProduksi
	for rows.Next() {
		var item model.FormulaProduksi
		if err := rows.Scan(&item.ID, &item.BarangJadi, &item.Kuantitas, &item.Satuan, &item.BahanBaku); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse formula produksi"})
			return
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result})
}

// GetFormulaProduksiByID retrieves a specific formulaProduksi by its ID
func GetFormulaProduksiByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `
		SELECT id, barang_jadi, kuantitas, satuan, bahan_baku
		FROM "formulaProduksi"
		WHERE id = $1
	`

	row := db.GetDB().QueryRow(query, id)
	var item model.FormulaProduksi
	err = row.Scan(&item.ID, &item.BarangJadi, &item.Kuantitas, &item.Satuan, &item.BahanBaku)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Formula Produksi not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch formula produksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    item})
}

// AddFormulaProduksi adds a new formulaProduksi record
func AddFormulaProduksi(c *gin.Context) {
	var input model.FormulaProduksi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `
		INSERT INTO "formulaProduksi" (barang_jadi, kuantitas, satuan, bahan_baku)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := db.GetDB().QueryRow(query, input.BarangJadi, input.Kuantitas, input.Satuan, input.BahanBaku).Scan(&input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add formula produksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Formula Produksi Added Successfully",
		"data":    input})
}

// UpdateFormulaProduksi dynamically updates formulaProduksi based on the provided fields
func UpdateFormulaProduksi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input model.FormulaProduksi
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Prepare the dynamic update query
	var updates []string
	var values []interface{}
	argPos := 1

	// Add fields only if they are provided
	if input.BarangJadi != "" {
		updates = append(updates, fmt.Sprintf("barang_jadi = $%d", argPos))
		values = append(values, input.BarangJadi)
		argPos++
	}
	if input.Kuantitas != 0 {
		updates = append(updates, fmt.Sprintf("kuantitas = $%d", argPos))
		values = append(values, input.Kuantitas)
		argPos++
	}
	if input.Satuan != 0 {
		updates = append(updates, fmt.Sprintf("satuan = $%d", argPos))
		values = append(values, input.Satuan)
		argPos++
	}
	if input.BahanBaku != "" {
		updates = append(updates, fmt.Sprintf("bahan_baku = $%d", argPos))
		values = append(values, input.BahanBaku)
		argPos++
	}

	// If no fields are provided, return an error
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Build the dynamic SQL query
	query := fmt.Sprintf(`
		UPDATE "formulaProduksi"
		SET %s
		WHERE id = $%d
	`, strings.Join(updates, ", "), argPos)

	// Add the ID to the end of the values array
	values = append(values, id)

	// Execute the update query
	_, err = db.GetDB().Exec(query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update formula produksi"})
		return
	}

	// Return the success message
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Formula Produksi Updated Successfully",
	})
}

// DeleteFormulaProduksi deletes a formulaProduksi record by its ID
func DeleteFormulaProduksi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `DELETE FROM "formulaProduksi" WHERE id = $1`
	_, err = db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete formula produksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Formula Produksi Deleted Successfully"})
}
