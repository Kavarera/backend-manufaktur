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

func ListMentah(c *gin.Context) {
	query := `
	SELECT bm."id", bm."nama", bm."kode_barang", bm."harga_standar",
	       bm."satuan", st."nama" AS satuan_nama,
		   bm."stok", bm."jumlah_turunan", bm."gudang", g."nama" AS gudang_nama,
		   bm."satuan_utama", su."nama" AS satuan_utama_nama 
	FROM "barangMentah" bm
	LEFT JOIN "satuanTurunan" st ON bm."satuan" = st."id"
	LEFT JOIN "gudang" g ON bm."gudang" = g."id"
	LEFT JOIN "barangSatuan" su ON bm."satuan_utama" = su."id"
	ORDER BY bm."kode_barang", bm."id"
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch barang list"})
		return
	}
	defer rows.Close()

	grouped := map[string]*model.BarangMentah{}

	for rows.Next() {
		var (
			id              int
			nama            string
			kodeBarang      string
			hargaStandar    float64
			satuanID        *int
			satuanNama      *string
			stok            float64
			gudangID        int
			gudangNama      string
			satuanUtamaID   *int
			satuanUtamaNama *string
			jumlahTurunan   sql.NullFloat64
		)

		err := rows.Scan(&id, &nama, &kodeBarang, &hargaStandar,
			&satuanID, &satuanNama, &stok, &jumlahTurunan, &gudangID, &gudangNama,
			&satuanUtamaID, &satuanUtamaNama)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to parse barang data" + err.Error()})
			return
		}

		if _, exists := grouped[kodeBarang]; !exists {
			grouped[kodeBarang] = &model.BarangMentah{
				ID:              id,
				Nama:            nama,
				KodeBarang:      kodeBarang,
				HargaStandar:    hargaStandar,
				GudangID:        gudangID,
				GudangNama:      gudangNama,
				SatuanUtamaID:   satuanUtamaID,
				SatuanUtamaNama: satuanUtamaNama,
				SatuanTurunan:   []model.BarangSatuanTurunanMentah{},
				Stok:            stok,
			}
		}

		if satuanID != nil && satuanNama != nil && jumlahTurunan.Valid {
			grouped[kodeBarang].SatuanTurunan = append(grouped[kodeBarang].SatuanTurunan, model.BarangSatuanTurunanMentah{
				SatuanID:      *satuanID,
				SatuanNama:    *satuanNama,
				JumlahTurunan: jumlahTurunan.Float64,
			})
		}
	}

	var result []model.BarangMentah
	for _, v := range grouped {
		result = append(result, *v)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}

func AddMentah(c *gin.Context) {
	var barang model.BarangMentah

	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid request payload"})
		return
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	var insertedIDs []int

	if len(barang.SatuanTurunan) > 0 {
		for _, st := range barang.SatuanTurunan {
			insertQuery := `
			INSERT INTO "barangMentah" 
			(nama, kode_barang, harga_standar, stok, gudang, satuan_utama, satuan, jumlah_turunan)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
			`
			var insertedID int
			err := tx.QueryRow(insertQuery,
				barang.Nama,
				barang.KodeBarang,
				barang.HargaStandar,
				barang.Stok,
				barang.GudangID,
				barang.SatuanUtamaID,
				st.SatuanID,
				st.JumlahTurunan,
			).Scan(&insertedID)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Insert failed: " + err.Error()})
				return
			}
			insertedIDs = append(insertedIDs, insertedID)
		}
	} else {
		// If no satuanTurunan
		satuanID := barang.SatuanID
		if satuanID == nil && barang.SatuanUtamaID != nil {
			satuanID = barang.SatuanUtamaID
		}
		if satuanID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Either satuanTurunan or satuan/satuanUtama is required"})
			return
		}

		insertQuery := `
		INSERT INTO "barangMentah" 
		(nama, kode_barang, harga_standar, stok, gudang, satuan_utama, satuan)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`
		var insertedID int
		err := tx.QueryRow(insertQuery,
			barang.Nama,
			barang.KodeBarang,
			barang.HargaStandar,
			barang.Stok,
			barang.GudangID,
			barang.SatuanUtamaID,
			*satuanID,
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

	barang.ID = insertedIDs[0]

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Barang mentah created successfully",
		"data":    barang,
	})
}

func UpdateMentah(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid id"})
		return
	}

	var barang model.BarangMentah

	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid request payload"})
		return
	}

	setClauses := []string{}
	values := []interface{}{}
	argPos := 1

	if barang.Nama != "" {
		setClauses = append(setClauses, fmt.Sprintf(`"nama"=$%d`, argPos))
		values = append(values, barang.Nama)
		argPos++
	}
	if barang.KodeBarang != "" {
		setClauses = append(setClauses, fmt.Sprintf(`"kode_barang"=$%d`, argPos))
		values = append(values, barang.KodeBarang)
		argPos++
	}
	if barang.HargaStandar != 0 {
		setClauses = append(setClauses, fmt.Sprintf(`"harga_standar"=$%d`, argPos))
		values = append(values, barang.HargaStandar)
		argPos++
	}
	if barang.Stok != 0 {
		setClauses = append(setClauses, fmt.Sprintf(`"stok"=$%d`, argPos))
		values = append(values, barang.Stok)
		argPos++
	}
	if barang.GudangID != 0 {
		setClauses = append(setClauses, fmt.Sprintf(`"gudang"=$%d`, argPos))
		values = append(values, barang.GudangID)
		argPos++
	}
	if barang.SatuanUtamaID != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"satuan_utama"=$%d`, argPos))
		values = append(values, barang.SatuanUtamaID)
		argPos++
	}
	if len(barang.SatuanTurunan) > 0 {
		turunan := barang.SatuanTurunan[0]

		if turunan.JumlahTurunan != 0 {
			setClauses = append(setClauses, fmt.Sprintf(`"jumlah_turunan"=$%d`, argPos))
			values = append(values, turunan.JumlahTurunan)
			argPos++
		}

		if turunan.SatuanID != 0 {
			setClauses = append(setClauses, fmt.Sprintf(`"satuan"=$%d`, argPos))
			values = append(values, turunan.SatuanID)
			argPos++
		}
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "No fields to update"})
		return
	}

	values = append(values, id)
	query := fmt.Sprintf(`UPDATE "barangMentah" SET %s WHERE "id"=$%d`, strings.Join(setClauses, ", "), argPos)

	res, err := db.GetDB().Exec(query, values...)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to update barang"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Barang updated successfully"})
}

func DeleteMentah(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid id"})
		return
	}

	query := `DELETE FROM "barangMentah" WHERE "id"=$1`
	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to delete barang"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Barang deleted successfully"})
}

func DeleteTurunanMentah(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid ID"})
		return
	}

	query := `
		UPDATE "barangMentah"
		SET "satuan" = NULL,
		    "jumlah_turunan" = NULL
		WHERE "id" = $1
	`

	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to clear satuan fields: " + err.Error()})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Satuan Turunan fields cleared successfully",
	})
}
