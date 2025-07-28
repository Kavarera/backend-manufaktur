package handler

import (
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListBarangProduksi(c *gin.Context) {
	query := `
	SELECT 
		bp.id, bp.nama, bp.kode_barang, bp.harga_standar, bp.harga_real,
		bp.satuan_utama, bs.nama, bp.stok, bp.gudang, g.nama,
		st.id, st.nama
	FROM "barangProduksi" bp
	JOIN "satuanTurunan" st ON bp.satuan = st.id
	JOIN "gudang" g ON bp.gudang = g.id
	LEFT JOIN "barangSatuan" bs ON bp.satuan_utama = bs.id
	ORDER BY bp.kode_barang, bp.id
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch barang produksi"})
		return
	}
	defer rows.Close()

	grouped := map[string]*model.BarangProduksi{}

	for rows.Next() {
		var (
			id              int
			nama            string
			kodeBarang      string
			hargaStandar    float64
			hargaReal       float64
			satuanUtamaID   *int
			satuanUtamaNama *string
			stok            float64
			gudangID        int
			gudangNama      string
			satuanID        int
			satuanNama      string
		)

		err := rows.Scan(&id, &nama, &kodeBarang, &hargaStandar, &hargaReal,
			&satuanUtamaID, &satuanUtamaNama, &stok, &gudangID, &gudangNama,
			&satuanID, &satuanNama)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to parse barang produksi"})
			return
		}

		if _, exists := grouped[kodeBarang]; !exists {
			grouped[kodeBarang] = &model.BarangProduksi{
				ID:              id,
				Nama:            nama,
				KodeBarang:      kodeBarang,
				HargaStandar:    hargaStandar,
				HargaReal:       hargaReal,
				SatuanUtamaID:   satuanUtamaID,
				SatuanUtamaNama: satuanUtamaNama,
				Stok:            stok,
				GudangID:        gudangID,
				GudangNama:      gudangNama,
			}
		}

		grouped[kodeBarang].SatuanTurunan = append(grouped[kodeBarang].SatuanTurunan, model.BarangSatuanTurunan{
			SatuanTurunanID:   satuanID,
			SatuanTurunanNama: satuanNama,
			Jumlah:            stok, // stok is jumlah for that satuanTurunan
		})
	}

	var result []model.BarangProduksi
	for _, v := range grouped {
		result = append(result, *v)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}

func GetBarangProduksiByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid ID"})
		return
	}

	var kodeBarang string
	err = db.GetDB().QueryRow(`SELECT kode_barang FROM "barangProduksi" WHERE id = $1`, id).Scan(&kodeBarang)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang produksi not found"})
		return
	}

	query := `
	SELECT 
		bp.id, bp.nama, bp.kode_barang, bp.harga_standar, bp.harga_real,
		bp.satuan_utama, bs.nama, bp.stok, bp.gudang, g.nama,
		st.id, st.nama
	FROM "barangProduksi" bp
	JOIN "satuanTurunan" st ON bp.satuan = st.id
	JOIN "gudang" g ON bp.gudang = g.id
	LEFT JOIN "barangSatuan" bs ON bp.satuan_utama = bs.id
	WHERE bp.kode_barang = $1
	`

	rows, err := db.GetDB().Query(query, kodeBarang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Query failed"})
		return
	}
	defer rows.Close()

	var result model.BarangProduksi
	var satuanTurunans []model.BarangSatuanTurunan

	for rows.Next() {
		var (
			id              int
			satuanUtamaID   *int
			satuanUtamaNama *string
			st              model.BarangSatuanTurunan
		)

		err := rows.Scan(&id, &result.Nama, &result.KodeBarang, &result.HargaStandar, &result.HargaReal,
			&satuanUtamaID, &satuanUtamaNama, &st.Jumlah, &result.GudangID, &result.GudangNama,
			&st.SatuanTurunanID, &st.SatuanTurunanNama)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Scan failed"})
			return
		}

		result.ID = id
		result.SatuanUtamaID = satuanUtamaID
		result.SatuanUtamaNama = satuanUtamaNama

		satuanTurunans = append(satuanTurunans, st)
	}

	result.SatuanTurunan = satuanTurunans

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}

func AddBarangProduksi(c *gin.Context) {
	var bp model.BarangProduksi

	if err := c.ShouldBindJSON(&bp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid request payload"})
		return
	}

	if len(bp.SatuanTurunan) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "At least one satuanTurunan is required"})
		return
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	insertQuery := `
		INSERT INTO "barangProduksi" 
		(nama, kode_barang, harga_standar, harga_real, satuan_utama, satuan, stok, gudang)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var insertedIDs []int

	for _, st := range bp.SatuanTurunan {
		var insertedID int
		err := tx.QueryRow(insertQuery,
			bp.Nama,
			bp.KodeBarang,
			bp.HargaStandar,
			bp.HargaReal,
			bp.SatuanUtamaID,
			st.SatuanTurunanID,
			st.Jumlah,
			bp.GudangID,
		).Scan(&insertedID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Insert failed: " + err.Error()})
			return
		}

		insertedIDs = append(insertedIDs, insertedID)
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Transaction commit failed"})
		return
	}

	// Respond with the same input plus a representative ID
	bp.ID = insertedIDs[0]

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Barang produksi created successfully",
		"data":    bp,
	})
}

func UpdateBarangProduksi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid id"})
		return
	}

	var bp model.BarangProduksi
	if err := c.ShouldBindJSON(&bp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid request payload"})
		return
	}

	// Validate that satuanTurunan is not empty
	if len(bp.SatuanTurunan) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "At least one satuanTurunan is required"})
		return
	}

	satuanID := bp.SatuanTurunan[0].SatuanTurunanID

	query := `
		UPDATE "barangProduksi"
		SET "nama"=$1, "kode_barang"=$2, "harga_standar"=$3, "harga_real"=$4, 
			"satuan"=$5, "stok"=$6, "gudang"=$7, "satuan_utama"=$8
		WHERE "id"=$9
	`

	res, err := db.GetDB().Exec(query, bp.Nama, bp.KodeBarang, bp.HargaStandar, bp.HargaReal,
		satuanID, bp.Stok, bp.GudangID, bp.SatuanUtamaID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to update barang produksi: " + err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang produksi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Barang produksi updated successfully",
	})
}

func DeleteBarangProduksi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid id"})
		return
	}

	query := `DELETE FROM "barangProduksi" WHERE "id"=$1`
	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to delete barang produksi"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang produksi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Barang produksi deleted successfully",
	})
}
