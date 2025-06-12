package handler

import (
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddPengambilanBarangBaku creates a new pengambilan barang baku and updates related tables
func AddPengambilanBarangBaku(c *gin.Context) {
	var input model.PengambilanBarangBaku
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Insert pengambilanBarangBaku
	query := `
		INSERT INTO "pengambilanBarangBaku" (id_perintah_kerja, id_barang_mentah, kebutuhan)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := db.GetDB().QueryRow(query, input.IDPerintahKerja, input.IDBarangMentah, input.Kebutuhan).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add pengambilan barang baku"})
		return
	}

	// Update hasil in perintahKerja table based on kebutuhan
	queryPerintahKerja := `
		UPDATE "perintahKerja"
		SET hasil = $1
		WHERE id = $2
	`
	_, err = db.GetDB().Exec(queryPerintahKerja, input.Kebutuhan, input.IDPerintahKerja)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hasil in perintahKerja"})
		return
	}

	// Update stok in barangMentah table
	queryBarangMentah := `
		UPDATE "barangMentah"
		SET stok = stok - $1
		WHERE id = $2
	`
	_, err = db.GetDB().Exec(queryBarangMentah, input.Kebutuhan, input.IDBarangMentah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stok in barangMentah"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengambilan Barang Baku added successfully",
		"data":    input,
	})
}

// GetPengambilanBarangBaku lists all pengambilanBarangBaku for a given perintahKerja ID
func GetPengambilanBarangBaku(c *gin.Context) {
	query := `
		SELECT 
			pbb.id,
			pbb.id_perintah_kerja,
			pbb.id_barang_mentah,
			pbb.kebutuhan,
			pk.tanggal_rilis,
			pk.tanggal_progres,
			pk.tanggal_selesai,
			pk.status AS perintah_kerja_status,
			bm.nama AS barang_mentah_nama,
			bm.kode_barang AS barang_mentah_kode,
			bm.harga_standar AS barang_mentah_harga_standar,
			bm.stok AS barang_mentah_stok
		FROM 
			"pengambilanBarangBaku" pbb
		JOIN 
			"perintahKerja" pk ON pk.id = pbb.id_perintah_kerja
		JOIN 
			"barangMentah" bm ON bm.id = pbb.id_barang_mentah
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		fmt.Println("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pengambilan barang baku"})
		return
	}
	defer rows.Close()

	var result []model.PengambilanBarangBaku
	for rows.Next() {
		var item model.PengambilanBarangBaku
		// Scan the data into the result object
		if err := rows.Scan(
			&item.ID,
			&item.IDPerintahKerja,
			&item.IDBarangMentah,
			&item.Kebutuhan,
			&item.TanggalRilis,
			&item.TanggalProgres,
			&item.TanggalSelesai,
			&item.StatusPerintahKerja,
			&item.NamaBarangMentah,
			&item.KodeBarangMentah,
			&item.HargaStandarBarangMentah,
			&item.StokBarangMentah,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse pengambilan barang baku"})
			return
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdatePengambilanBarangBaku updates pengambilanBarangBaku and adjusts related fields
func UpdatePengambilanBarangBaku(c *gin.Context) {
	id := c.Param("id")

	var input model.PengambilanBarangBaku
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Get current values
	var oldKebutuhan float64
	var idPerintahKerja string
	var idBarangMentah int
	query := `SELECT id_perintah_kerja, id_barang_mentah, kebutuhan FROM "pengambilanBarangBaku" WHERE id = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&idPerintahKerja, &idBarangMentah, &oldKebutuhan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing record"})
		return
	}

	// Update pengambilanBarangBaku record
	queryUpdate := `
		UPDATE "pengambilanBarangBaku"
		SET id_barang_mentah = $1, kebutuhan = $2
		WHERE id = $3
	`
	_, err = db.GetDB().Exec(queryUpdate, input.IDBarangMentah, input.Kebutuhan, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pengambilan barang baku"})
		return
	}

	// Adjust "hasil" in perintahKerja
	queryPerintahKerja := `
		UPDATE "perintahKerja"
		SET hasil = $1 
		WHERE id = $2
	`
	_, err = db.GetDB().Exec(queryPerintahKerja, input.Kebutuhan, idPerintahKerja)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hasil in perintahKerja"})
		return
	}

	// Adjust "stok" in barangMentah
	queryBarangMentah := `
		UPDATE "barangMentah"
		SET stok = stok - $1 + $2
		WHERE id = $3
	`
	_, err = db.GetDB().Exec(queryBarangMentah, input.Kebutuhan, oldKebutuhan, idBarangMentah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stok in barangMentah"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengambilan Barang Baku updated successfully"})
}

// DeletePengambilanBarangBaku deletes a specific pengambilanBarangBaku and reverts related fields
func DeletePengambilanBarangBaku(c *gin.Context) {
	id := c.Param("id")

	// Get the current values
	var kebutuhan float64
	var idPerintahKerja string
	var idBarangMentah int
	query := `SELECT id_perintah_kerja, id_barang_mentah, kebutuhan FROM "pengambilanBarangBaku" WHERE id = $1`
	err := db.GetDB().QueryRow(query, id).Scan(&idPerintahKerja, &idBarangMentah, &kebutuhan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing record"})
		return
	}

	// Delete pengambilanBarangBaku record
	queryDelete := `DELETE FROM "pengambilanBarangBaku" WHERE id = $1`
	_, err = db.GetDB().Exec(queryDelete, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pengambilan barang baku"})
		return
	}

	// Revert "hasil" in perintahKerja
	queryPerintahKerja := `
		UPDATE "perintahKerja"
		SET hasil = hasil - $1
		WHERE id = $2
	`
	_, err = db.GetDB().Exec(queryPerintahKerja, kebutuhan, idPerintahKerja)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revert hasil in perintahKerja"})
		return
	}

	// Revert "stok" in barangMentah
	queryBarangMentah := `
		UPDATE "barangMentah"
		SET stok = stok + $1
		WHERE id = $2
	`
	_, err = db.GetDB().Exec(queryBarangMentah, kebutuhan, idBarangMentah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revert stok in barangMentah"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengambilan Barang Baku deleted successfully"})
}
