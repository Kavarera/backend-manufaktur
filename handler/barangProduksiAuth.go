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

func ListBarangProduksi(c *gin.Context) {
	query := `
	SELECT bp."id", bp."nama", bp."kode_barang", bp."harga_standar", bp."harga_real", 
	       bp."satuan", st."nama", bp."stok", bp."gudang", g."nama",
	       bp."satuan_utama", bs."nama"
	FROM "barangProduksi" bp
	JOIN "satuanTurunan" st ON bp."satuan" = st."id"
	JOIN "gudang" g ON bp."gudang" = g."id"
	LEFT JOIN "barangSatuan" bs ON bp."satuan_utama" = bs."id"
	ORDER BY bp."id";

	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch barang produksi"})
		return
	}
	defer rows.Close()

	var list []model.BarangProduksi
	for rows.Next() {
		var bp model.BarangProduksi
		err := rows.Scan(&bp.ID, &bp.Nama, &bp.KodeBarang, &bp.HargaStandar, &bp.HargaReal,
			&bp.SatuanID, &bp.SatuanNama, &bp.Stok, &bp.GudangID, &bp.GudangNama,
			&bp.SatuanUtamaID, &bp.SatuanUtamaNama)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to parse barang produksi"})
			return
		}
		list = append(list, bp)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    list,
	})
}

func GetBarangProduksiByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid id"})
		return
	}

	query := `
	SELECT bp."id", bp."nama", bp."kode_barang", bp."harga_standar", bp."harga_real", 
	       bp."satuan", st."nama", bp."stok", bp."gudang", g."nama",
	       bp."satuan_utama", bs."nama"
	FROM "barangProduksi" bp
	JOIN "satuanTurunan" st ON bp."satuan" = st."id"
	JOIN "gudang" g ON bp."gudang" = g."id"
	LEFT JOIN "barangSatuan" bs ON bp."satuan_utama" = bs."id"
	WHERE bp."id" = $1
	`

	row := db.GetDB().QueryRow(query, id)
	var bp model.BarangProduksi
	err = row.Scan(&bp.ID, &bp.Nama, &bp.KodeBarang, &bp.HargaStandar, &bp.HargaReal,
		&bp.SatuanID, &bp.SatuanNama, &bp.Stok, &bp.GudangID, &bp.GudangNama,
		&bp.SatuanUtamaID, &bp.SatuanUtamaNama)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang produksi not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch barang produksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    bp,
	})
}

func AddBarangProduksi(c *gin.Context) {
	var bp model.BarangProduksi

	if err := c.ShouldBindJSON(&bp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid request payload"})
		return
	}

	// Build dynamic query parts to allow optional satuanUtama
	columns := []string{"nama", "kode_barang", "harga_standar", "harga_real", "satuan", "stok", "gudang"}
	values := []interface{}{bp.Nama, bp.KodeBarang, bp.HargaStandar, bp.HargaReal, bp.SatuanID, bp.Stok, bp.GudangID}
	placeholders := []string{"$1", "$2", "$3", "$4", "$5", "$6", "$7"}
	argPos := 8

	if bp.SatuanUtamaID != nil {
		columns = append(columns, "satuan_utama")
		values = append(values, *bp.SatuanUtamaID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
		argPos++
	}

	query := fmt.Sprintf(`
	INSERT INTO "barangProduksi" (%s)
	VALUES (%s)
	RETURNING "id"
	`, strings.Join(columns, ","), strings.Join(placeholders, ","))

	err := db.GetDB().QueryRow(query, values...).Scan(&bp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to create barang produksi: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Berhasil",
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

	query := `
	UPDATE "barangProduksi"
	SET "nama"=$1, "kode_barang"=$2, "harga_standar"=$3, "harga_real"=$4, "satuan"=$5, "stok"=$6, "gudang"=$7, "satuan_utama"=$8
	WHERE "id"=$9
	`

	res, err := db.GetDB().Exec(query, bp.Nama, bp.KodeBarang, bp.HargaStandar, bp.HargaReal,
		bp.SatuanID, bp.Stok, bp.GudangID, bp.SatuanUtamaID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to update barang produksi"})
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
