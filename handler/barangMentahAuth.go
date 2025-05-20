package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListMentah(c *gin.Context) {
	query :=
		`
		SELECT bm."id", bm."nama", bm."kode_barang", bm."harga_standar", bm."satuan", st."nama" AS satuan_nama,
      		   bm."stok", bm."gudang", g."nama" AS gudang_nama
		FROM "barangMentah" bm
		LEFT JOIN "satuanTurunan" st ON bm."satuan" = st."id"
		LEFT JOIN "gudang" g ON bm."gudang" = g."id"
		ORDER BY bm."id";
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch barang list"})
		return
	}
	defer rows.Close()

	var barangList []model.BarangMentah

	for rows.Next() {
		var barang model.BarangMentah
		err := rows.Scan(&barang.ID, &barang.Nama, &barang.KodeBarang, &barang.HargaStandar,
			&barang.SatuanID, &barang.SatuanNama, &barang.Stok, &barang.GudangID, &barang.GudangNama)
		if err != nil {
			if err != sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Barang Mentah Not Found"})
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse barang data"})
			return
		}
		barangList = append(barangList, barang)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    barangList,
	})
}

func AddMentah(c *gin.Context) {
	var barang model.BarangMentah

	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid request payload"})
		return
	}

	query := `
	INSERT INTO "barangMentah" ("nama", "kode_barang", "harga_standar", "satuan", "stok", "gudang")
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING "id"
	`

	err := db.GetDB().QueryRow(query, barang.Nama, barang.KodeBarang, barang.HargaStandar,
		barang.SatuanID, barang.Stok, barang.GudangID).Scan(&barang.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to create barang"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Berhasil",
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

	query := `
	UPDATE "barangMentah"
	SET "nama"=$1, "kode_barang"=$2, "harga_standar"=$3, "satuan"=$4, "stok"=$5, "gudang"=$6
	WHERE "id"=$7
	`

	res, err := db.GetDB().Exec(query, barang.Nama, barang.KodeBarang, barang.HargaStandar,
		barang.SatuanID, barang.Stok, barang.GudangID, id)
	if err != nil {
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
