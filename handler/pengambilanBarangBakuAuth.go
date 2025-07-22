package handler

import (
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AddPengambilanBarangBaku creates a new pengambilan barang baku and updates related tables
func AddPengambilanBarangBaku(c *gin.Context) {
	var input struct {
		IDPerintahKerja string                        `json:"idPerintahKerja"`
		BarangBaku      []model.PengambilanBarangBaku `json:"barangBaku"`
	}

	// Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Get current timestamp for tanggal_waktu
	tanggalWaktu := time.Now()

	// Loop through each barangBaku in the payload and insert them
	for _, barang := range input.BarangBaku {
		// Insert pengambilanBarangBaku
		query := `
			INSERT INTO "pengambilanBarangBaku" (id_perintah_kerja, id_barang_mentah, kebutuhan, tanggal_waktu)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		var id int
		err := db.GetDB().QueryRow(query, input.IDPerintahKerja, barang.IDBarangMentah, barang.Kebutuhan, tanggalWaktu).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add pengambilan barang baku"})
			return
		}

		// Update stok in barangMentah table
		queryBarangMentah := `
			UPDATE "barangMentah"
			SET stok = stok - $1
			WHERE id = $2
		`
		_, err = db.GetDB().Exec(queryBarangMentah, barang.Kebutuhan, barang.IDBarangMentah)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stok in barangMentah"})
			return
		}
	}

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengambilan Barang Baku added successfully",
		"status":  "OK",
		"data":    input,
	})
}

func GetPengambilanBarangBaku(c *gin.Context) {
	query := `
		SELECT 
			pbb.id,
			pbb.id_perintah_kerja,
			pbb.id_barang_mentah,
			pbb.kebutuhan,
			pbb.tanggal_waktu,
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
		ORDER BY 
			pbb.id_perintah_kerja, pbb.id
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pengambilan barang baku"})
		return
	}
	defer rows.Close()

	// BarangMentah structure (a subset of your full model)
	type BarangMentah struct {
		ID                       int       `json:"idPengambilanBarangBaku"`
		IDBarangMentah           int       `json:"idBarangMentah"`
		NamaBarangMentah         string    `json:"namaBarangMentah"`
		KodeBarangMentah         string    `json:"kodeBarangMentah"`
		HargaStandarBarangMentah float64   `json:"hargaStandarBarangMentah"`
		StokBarangMentah         float64   `json:"stokBarangMentah"`
		Kebutuhan                float64   `json:"kebutuhan"`
		TanggalWaktu             time.Time `json:"tanggalWaktu"`
	}

	// Grouped structure per Perintah Kerja
	type GroupedPengambilan struct {
		IDPerintahKerja     string         `json:"idPerintahKerja"`
		TanggalRilis        time.Time      `json:"tanggalRilis"`
		TanggalProgres      *time.Time     `json:"tanggalProgres"`
		TanggalSelesai      *time.Time     `json:"tanggalSelesai"`
		StatusPerintahKerja string         `json:"statusPerintahKerja"`
		BarangMentah        []BarangMentah `json:"barangMentah"`
	}

	grouped := make(map[string]*GroupedPengambilan)

	for rows.Next() {
		var record model.PengambilanBarangBaku
		err := rows.Scan(
			&record.ID,
			&record.IDPerintahKerja,
			&record.IDBarangMentah,
			&record.Kebutuhan,
			&record.TanggalWaktu,
			&record.TanggalRilis,
			&record.TanggalProgres,
			&record.TanggalSelesai,
			&record.StatusPerintahKerja,
			&record.NamaBarangMentah,
			&record.KodeBarangMentah,
			&record.HargaStandarBarangMentah,
			&record.StokBarangMentah,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse result"})
			return
		}

		group, exists := grouped[record.IDPerintahKerja]
		if !exists {
			group = &GroupedPengambilan{
				IDPerintahKerja:     record.IDPerintahKerja,
				TanggalRilis:        record.TanggalRilis,
				TanggalProgres:      record.TanggalProgres,
				TanggalSelesai:      record.TanggalSelesai,
				StatusPerintahKerja: record.StatusPerintahKerja,
				BarangMentah:        []BarangMentah{},
			}
			grouped[record.IDPerintahKerja] = group
		}

		group.BarangMentah = append(group.BarangMentah, BarangMentah{
			ID:                       record.ID,
			IDBarangMentah:           record.IDBarangMentah,
			NamaBarangMentah:         record.NamaBarangMentah,
			KodeBarangMentah:         record.KodeBarangMentah,
			HargaStandarBarangMentah: record.HargaStandarBarangMentah,
			StokBarangMentah:         record.StokBarangMentah,
			Kebutuhan:                record.Kebutuhan,
			TanggalWaktu:             record.TanggalWaktu,
		})
	}

	// Convert map to slice
	var response []GroupedPengambilan
	for _, g := range grouped {
		response = append(response, *g)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    response,
	})
}

func UpdatePengambilanBarangBaku(c *gin.Context) {
	idPerintahKerja := c.Param("idPerintahKerja")

	var input struct {
		BarangBaku []struct {
			ID        int     `json:"idPengambilanBarangBaku"`
			Kebutuhan float64 `json:"kebutuhan"`
		} `json:"barangBaku"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if len(input.BarangBaku) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided"})
		return
	}

	dbConn := db.GetDB()
	tx, err := dbConn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, barang := range input.BarangBaku {
		var oldKebutuhan float64
		var idBarangMentah int

		// 1. Fetch current record
		queryGet := `
			SELECT id_barang_mentah, kebutuhan
			FROM "pengambilanBarangBaku"
			WHERE id = $1 AND id_perintah_kerja = $2
		`
		err = tx.QueryRow(queryGet, barang.ID, idPerintahKerja).Scan(&idBarangMentah, &oldKebutuhan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch record for ID %d", barang.ID)})
			return
		}

		// 2. Update pengambilanBarangBaku
		queryUpdate := `
			UPDATE "pengambilanBarangBaku"
			SET kebutuhan = $1, tanggal_waktu = $2
			WHERE id = $3
		`
		_, err = tx.Exec(queryUpdate, barang.Kebutuhan, time.Now(), barang.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update record for ID %d", barang.ID)})
			return
		}

		// 3. Revert old kebutuhan from stok
		queryRevertStok := `
			UPDATE "barangMentah"
			SET stok = stok + $1
			WHERE id = $2
		`
		_, err = tx.Exec(queryRevertStok, oldKebutuhan, idBarangMentah)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to revert stok for ID %d", barang.ID)})
			return
		}

		// 4. Subtract new kebutuhan from stok
		querySubtractStok := `
			UPDATE "barangMentah"
			SET stok = stok - $1
			WHERE id = $2
		`
		_, err = tx.Exec(querySubtractStok, barang.Kebutuhan, idBarangMentah)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update stok for ID %d", barang.ID)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Pengambilan Barang Baku updated successfully",
	})
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

	c.JSON(http.StatusOK, gin.H{"message": "Pengambilan Barang Baku deleted successfully", "status": "OK"})
}
