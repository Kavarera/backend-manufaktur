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

func ListFormulaProduksi(c *gin.Context) {
	query := `
		SELECT id, barang_jadi, kuantitas, satuan, bahan_baku, satuan_turunan
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
		var satuan sql.NullFloat64
		var satuanTurunan sql.NullFloat64

		err := rows.Scan(
			&item.ID,
			&item.BarangJadi,
			&item.Kuantitas,
			&satuan,
			&item.BahanBaku,
			&satuanTurunan,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse formula produksi"})
			return
		}

		if satuan.Valid {
			item.Satuan = satuan.Float64
		} else {
			item.Satuan = 0
		}

		if satuanTurunan.Valid {
			item.SatuanTurunan = satuanTurunan.Float64
		} else {
			item.SatuanTurunan = 0
		}

		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}

func GetFormulaProduksiByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `
		SELECT id, barang_jadi, kuantitas, satuan, bahan_baku, satuan_turunan
		FROM "formulaProduksi"
		WHERE id = $1
	`

	row := db.GetDB().QueryRow(query, id)

	var item model.FormulaProduksi
	var satuan sql.NullFloat64
	var satuanTurunan sql.NullFloat64

	err = row.Scan(
		&item.ID,
		&item.BarangJadi,
		&item.Kuantitas,
		&satuan,
		&item.BahanBaku,
		&satuanTurunan,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Formula Produksi not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch formula produksi"})
		return
	}

	if satuan.Valid {
		item.Satuan = satuan.Float64
	}
	if satuanTurunan.Valid {
		item.SatuanTurunan = satuanTurunan.Float64
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    item,
	})
}

// AddFormulaProduksi adds one or more new formulaProduksi records
func AddFormulaProduksi(c *gin.Context) {
	var inputs []model.FormulaProduksi
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `
		INSERT INTO "formulaProduksi" 
		(barang_jadi, kuantitas, satuan, bahan_baku, satuan_turunan)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	dbConn := db.GetDB()
	tx, err := dbConn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var added []model.FormulaProduksi

	for i, input := range inputs {
		var id int
		err := tx.QueryRow(
			query,
			input.BarangJadi,
			input.Kuantitas,
			sql.NullFloat64{Float64: input.Satuan, Valid: input.Satuan != 0},
			input.BahanBaku,
			sql.NullFloat64{Float64: input.SatuanTurunan, Valid: input.SatuanTurunan != 0},
		).Scan(&id)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed at item %d: %v", i+1, err)})
			return
		}
		input.ID = id
		added = append(added, input)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Formula Produksi Added Successfully",
		"data":    added,
	})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var updates []string
	var values []interface{}
	argPos := 1

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
	if input.SatuanTurunan != 0 {
		updates = append(updates, fmt.Sprintf("satuan_turunan = $%d", argPos))
		values = append(values, input.SatuanTurunan)
		argPos++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	query := fmt.Sprintf(`
		UPDATE "formulaProduksi"
		SET %s
		WHERE id = $%d
	`, strings.Join(updates, ", "), argPos)

	values = append(values, id)

	_, err = db.GetDB().Exec(query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update formula produksi"})
		return
	}

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
