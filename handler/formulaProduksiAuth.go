package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ListFormulaProduksi(c *gin.Context) {
	query := `
		SELECT id, id_barang_produksi, kuantitas, tanggal_mulai, nama_produksi
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
		var tanggalMulai time.Time
		err := rows.Scan(&item.ID, &item.IDBarangProduksi, &item.Kuantitas, &tanggalMulai, &item.NamaProduksi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse formula produksi"})
			return
		}
		item.TanggalMulai = tanggalMulai
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}

func GetFormulaProduksiByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `
		SELECT id, id_barang_produksi, kuantitas, tanggal_mulai, nama_produksi
		FROM "formulaProduksi"
		WHERE id = $1
	`
	var item model.FormulaProduksi
	var tanggalMulai time.Time
	err = db.GetDB().QueryRow(query, id).Scan(&item.ID, &item.IDBarangProduksi, &item.Kuantitas, &tanggalMulai, &item.NamaProduksi)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Formula Produksi not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch formula produksi"})
		return
	}
	item.TanggalMulai = tanggalMulai

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Berhasil", "data": item})
}

func AddFormulaProduksi(c *gin.Context) {
	var inputs []model.FormulaProduksi
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `
		INSERT INTO "formulaProduksi" (id_barang_produksi, kuantitas, tanggal_mulai, nama_produksi)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	tx, err := db.GetDB().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var added []model.FormulaProduksi
	for i, input := range inputs {
		var id int
		err := tx.QueryRow(query, input.IDBarangProduksi, input.Kuantitas, input.TanggalMulai, input.NamaProduksi).Scan(&id)
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

func UpdateFormulaProduksi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	allowedFields := map[string]bool{
		"idBarangProduksi": true,
		"kuantitas":        true,
		"tanggalMulai":     true,
		"namaProduksi":     true,
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	for key, value := range input {
		if !allowedFields[key] {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Field %s not allowed for update", key)})
			return
		}

		column := map[string]string{
			"idBarangProduksi": "id_barang_produksi",
			"kuantitas":        "kuantitas",
			"tanggalMulai":     "tanggal_mulai",
			"namaProduksi":     "nama_produksi",
		}[key]

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIdx))
		args = append(args, value)
		argIdx++
	}

	query := fmt.Sprintf(`UPDATE "formulaProduksi" SET %s WHERE id = $%d`,
		strings.Join(setClauses, ", "), argIdx)

	args = append(args, id)

	_, err = db.GetDB().Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update formula produksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Formula Produksi Updated Successfully"})
}

func DeleteFormulaProduksi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Formula Produksi Deleted Successfully"})
}
