package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ListRencanaProduksi lists all production plans
func ListRencanaProduksi(c *gin.Context) {
	query := `SELECT "id", "id_barang_produksi", "tanggal_mulai", "tanggal_selesai" FROM "rencanaProduksi" ORDER BY "id"`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rencana produksi"})
		return
	}
	defer rows.Close()

	var list []model.RencanaProduksi
	for rows.Next() {
		var rp model.RencanaProduksi
		err := rows.Scan(&rp.ID, &rp.BarangProduksiID, &rp.TanggalMulai, &rp.TanggalSelesai)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rencana produksi"})
			return
		}
		list = append(list, rp)
	}

	c.JSON(http.StatusOK, list)
}

// GetRencanaProduksiByID gets a production plan by ID
func GetRencanaProduksiByID(c *gin.Context) {
	id := c.Param("id")

	query := `SELECT "id", "id_barang_produksi", "tanggal_mulai", "tanggal_selesai" FROM "rencanaProduksi" WHERE "id"=$1`
	row := db.GetDB().QueryRow(query, id)

	var rp model.RencanaProduksi
	err := row.Scan(&rp.ID, &rp.BarangProduksiID, &rp.TanggalMulai, &rp.TanggalSelesai)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rencana produksi not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rencana produksi"})
		}
		return
	}

	c.JSON(http.StatusOK, rp)
}

// AddRencanaProduksi creates a new production plan
func AddRencanaProduksi(c *gin.Context) {
	var payload model.RencanaProduksiAdd

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Combine tanggalMulai + waktuMulai
	startStr := fmt.Sprintf("%sT%s", payload.TanggalMulai, payload.WaktuMulai)
	tanggalMulai, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggalMulai or waktuMulai"})
		return
	}

	// Combine tanggalSelesai + waktuSelesai
	endStr := fmt.Sprintf("%sT%s", payload.TanggalSelesai, payload.WaktuSelesai)
	tanggalSelesai, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tanggalSelesai or waktuSelesai"})
		return
	}

	query := `
	INSERT INTO "rencanaProduksi" ("id", "id_barang_produksi", "tanggal_mulai", "tanggal_selesai")
	VALUES ($1, $2, $3, $4)
	`

	_, err = db.GetDB().Exec(query, payload.ID, payload.BarangProduksiID, tanggalMulai, tanggalSelesai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rencana produksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               payload.ID,
		"barangProduksiId": payload.BarangProduksiID,
		"tanggalMulai":     tanggalMulai,
		"tanggalSelesai":   tanggalSelesai,
	})
}

// UpdateRencanaProduksi updates a production plan by ID
func UpdateRencanaProduksi(c *gin.Context) {
	id := c.Param("id")

	var payload model.RencanaProduksi
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var setClauses []string
	var args []interface{}
	argPos := 1

	if payload.BarangProduksiID != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"id_barang_produksi"=$%d`, argPos))
		args = append(args, *payload.BarangProduksiID)
		argPos++
	}
	if payload.TanggalMulai != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_mulai"=$%d`, argPos))
		args = append(args, *payload.TanggalMulai)
		argPos++
	}
	if payload.TanggalSelesai != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_selesai"=$%d`, argPos))
		args = append(args, *payload.TanggalSelesai)
		argPos++
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	args = append(args, id)
	sql := fmt.Sprintf(`UPDATE "rencanaProduksi" SET %s WHERE "id"=$%d`, strings.Join(setClauses, ", "), argPos)

	res, err := db.GetDB().Exec(sql, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rencana produksi"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rencana produksi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rencana produksi updated successfully"})
}

// DeleteRencanaProduksi deletes a production plan by ID
func DeleteRencanaProduksi(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM "rencanaProduksi" WHERE "id"=$1`

	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rencana produksi"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rencana produksi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rencana produksi deleted successfully"})
}
