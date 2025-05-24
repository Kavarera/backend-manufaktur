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
	SELECT bm."id", bm."nama", bm."kode_barang", bm."harga_standar", bm."satuan", st."nama" AS satuan_nama,
			bm."stok", bm."gudang", g."nama" AS gudang_nama,
			bm."satuan_utama", su."nama" AS satuan_utama_nama
	FROM "barangMentah" bm
	LEFT JOIN "satuanTurunan" st ON bm."satuan" = st."id"
	LEFT JOIN "gudang" g ON bm."gudang" = g."id"
	LEFT JOIN "barangSatuan" su ON bm."satuan_utama" = su."id"
	ORDER BY bm."id";
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch barang list"})
		return
	}
	defer rows.Close()

	var barangList []model.BarangMentah

	for rows.Next() {
		var barang model.BarangMentah
		err := rows.Scan(&barang.ID, &barang.Nama, &barang.KodeBarang, &barang.HargaStandar,
			&barang.SatuanID, &barang.SatuanNama, &barang.Stok, &barang.GudangID, &barang.GudangNama,
			&barang.SatuanUtamaID, &barang.SatuanUtamaNama)
		if err != nil {
			if err != sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Barang Mentah Not Found"})
			}
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to parse barang data"})
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

	columns := []string{"nama", "kode_barang", "harga_standar", "stok", "gudang"}
	values := []interface{}{barang.Nama, barang.KodeBarang, barang.HargaStandar, barang.Stok, barang.GudangID}
	placeholders := []string{"$1", "$2", "$3", "$4", "$5"}
	argPos := 6

	if barang.SatuanID != nil {
		columns = append(columns, "satuan")
		values = append(values, *barang.SatuanID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
		argPos++
	}
	if barang.SatuanUtamaID != nil {
		columns = append(columns, "satuan_utama")
		values = append(values, *barang.SatuanUtamaID)
		placeholders = append(placeholders, fmt.Sprintf("$%d", argPos))
		argPos++
	}

	query := fmt.Sprintf(`
		INSERT INTO "barangMentah" (%s)
		VALUES (%s)
		RETURNING "id"
	`, strings.Join(columns, ","), strings.Join(placeholders, ","))

	err := db.GetDB().QueryRow(query, values...).Scan(&barang.ID) // scan id here
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to create barang: " + err.Error()})
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
	if barang.SatuanID != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"satuan"=$%d`, argPos))
		values = append(values, barang.SatuanID)
		argPos++
	}
	if barang.SatuanUtamaID != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"satuan_utama"=$%d`, argPos))
		values = append(values, barang.SatuanUtamaID)
		argPos++
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "No fields to update"})
		return
	}

	values = append(values, id)
	query := fmt.Sprintf(`UPDATE "barangMentah" SET %s WHERE "id"=$%d`, strings.Join(setClauses, ", "), argPos)

	res, err := db.GetDB().Exec(query, values...)
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
