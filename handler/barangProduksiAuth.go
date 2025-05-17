package handler

import (
	"database/sql"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListBarangProduksi lists all barangProduksi with related satuan and gudang names
func ListBarangProduksi(c *gin.Context) {
	query := `
	SELECT bp."id", bp."nama", bp."kode_barang", bp."harga_standar", bp."harga_real", 
	       bp."satuan", st."nama", bp."stok", bp."gudang", g."nama"
	FROM "barangProduksi" bp
	JOIN "satuanTurunan" st ON bp."satuan" = st."id"
	JOIN "gudang" g ON bp."gudang" = g."id"
	ORDER BY bp."id"
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch barang produksi"})
		return
	}
	defer rows.Close()

	var list []model.BarangProduksi
	for rows.Next() {
		var bp model.BarangProduksi
		err := rows.Scan(&bp.ID, &bp.Nama, &bp.KodeBarang, &bp.HargaStandar, &bp.HargaReal,
			&bp.SatuanID, &bp.SatuanNama, &bp.Stok, &bp.GudangID, &bp.GudangNama)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse barang produksi"})
			return
		}
		list = append(list, bp)
	}

	c.JSON(http.StatusOK, list)
}

// GetBarangProduksiByID returns single barangProduksi by ID
func GetBarangProduksiByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	query := `
	SELECT bp."id", bp."nama", bp."kode_barang", bp."harga_standar", bp."harga_real", 
	       bp."satuan", st."nama", bp."stok", bp."gudang", g."nama"
	FROM "barangProduksi" bp
	JOIN "satuanTurunan" st ON bp."satuan" = st."id"
	JOIN "gudang" g ON bp."gudang" = g."id"
	WHERE bp."id" = $1
	`

	row := db.GetDB().QueryRow(query, id)
	var bp model.BarangProduksi
	err = row.Scan(&bp.ID, &bp.Nama, &bp.KodeBarang, &bp.HargaStandar, &bp.HargaReal,
		&bp.SatuanID, &bp.SatuanNama, &bp.Stok, &bp.GudangID, &bp.GudangNama)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Barang produksi not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch barang produksi"})
		}
		return
	}

	c.JSON(http.StatusOK, bp)
}

// AddBarangProduksi creates a new barangProduksi
func AddBarangProduksi(c *gin.Context) {
	var bp model.BarangProduksi

	if err := c.ShouldBindJSON(&bp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `
	INSERT INTO "barangProduksi" ("nama", "kode_barang", "harga_standar", "harga_real", "satuan", "stok", "gudang")
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING "id"
	`

	err := db.GetDB().QueryRow(query, bp.Nama, bp.KodeBarang, bp.HargaStandar, bp.HargaReal,
		bp.SatuanID, bp.Stok, bp.GudangID).Scan(&bp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create barang produksi"})
		return
	}

	c.JSON(http.StatusCreated, bp)
}

// UpdateBarangProduksi updates barangProduksi by ID
func UpdateBarangProduksi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var bp model.BarangProduksi
	if err := c.ShouldBindJSON(&bp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `
	UPDATE "barangProduksi"
	SET "nama"=$1, "kode_barang"=$2, "harga_standar"=$3, "harga_real"=$4, "satuan"=$5, "stok"=$6, "gudang"=$7
	WHERE "id"=$8
	`

	res, err := db.GetDB().Exec(query, bp.Nama, bp.KodeBarang, bp.HargaStandar, bp.HargaReal,
		bp.SatuanID, bp.Stok, bp.GudangID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update barang produksi"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang produksi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang produksi updated successfully"})
}

// DeleteBarangProduksi deletes barangProduksi by ID
func DeleteBarangProduksi(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	query := `DELETE FROM "barangProduksi" WHERE "id"=$1`
	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete barang produksi"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang produksi not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang produksi deleted successfully"})
}
