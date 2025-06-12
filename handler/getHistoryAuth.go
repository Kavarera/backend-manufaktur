package handler

import (
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPerintahKerjaDetailsByID fetches Perintah Kerja along with related PenyelesaianBarangJadi, PengambilanBarangBaku, and BarangMentah data
func GetPerintahKerjaDetailsByID(c *gin.Context) {
	id := c.Param("id") // Get the perintahKerja ID from the URL

	// SQL query to join perintahKerja with related tables
	query := `
		SELECT 
			pk.id,
			pk.tanggal_rilis,
			pk.tanggal_progres,
			pk.tanggal_selesai,
			pk.status,
			pk.hasil,
			pk.customer,
			pk.keterangan,
			pk.document_url,
			pk.document_nama,
			pbd.id AS penyelesaian_id, 
			pbd.nama AS penyelesaian_nama,
			pbd.jumlah AS penyelesaian_jumlah,
			pbd.tanggal_penyelesaian,
			pbb.id AS pengambilan_id,
			pbb.kebutuhan AS pengambilan_kebutuhan,
			bm.id AS barang_mentah_id,
			bm.nama AS barang_mentah_nama,
			bm.kode_barang AS barang_mentah_kode,
			bm.harga_standar AS barang_mentah_harga_standar,
			bm.stok AS barang_mentah_stok
		FROM 
			"perintahKerja" pk
		LEFT JOIN 
			"penyelesaianBarangJadi" pbd ON pbd.id_perintah_kerja = pk.id
		LEFT JOIN 
			"pengambilanBarangBaku" pbb ON pbb.id_perintah_kerja = pk.id
		LEFT JOIN 
			"barangMentah" bm ON bm.id = pbb.id_barang_mentah
		WHERE pk.id = $1
	`

	rows, err := db.GetDB().Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch perintah kerja details"})
		return
	}
	defer rows.Close()

	// Initialize a PerintahKerjaDetails object
	var details model.PerintahKerjaDetails
	var penyelesaian []model.PenyelesaianBarangJadi
	var pengambilan []model.PengambilanBarangBaku
	var barangMentah []model.BarangMentah

	// Loop through the rows and scan the data
	for rows.Next() {
		var pk model.PerintahKerja
		var pbd model.PenyelesaianBarangJadi
		var pbb model.PengambilanBarangBaku
		var bm model.BarangMentah

		err := rows.Scan(
			&pk.ID, &pk.TanggalRilis, &pk.TanggalProgres, &pk.TanggalSelesai, &pk.Status,
			&pk.Hasil, &pk.Customer, &pk.Keterangan, &pk.DocumentURL, &pk.DocumentNama,
			&pbd.ID, &pbd.Nama, &pbd.Jumlah, &pbd.TanggalPenyelesaian,
			&pbb.ID, &pbb.Kebutuhan,
			&bm.ID, &bm.Nama, &bm.KodeBarang, &bm.HargaStandar, &bm.Stok,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}

		// Add the data to the appropriate slices
		if pbd.ID != 0 {
			penyelesaian = append(penyelesaian, pbd)
		}
		if pbb.ID != 0 {
			pengambilan = append(pengambilan, pbb)
		}
		if bm.ID != 0 {
			barangMentah = append(barangMentah, bm)
		}

		// Setting the PerintahKerja data
		details.PerintahKerja = pk
	}

	// Assign the results to the response struct
	details.PenyelesaianBarangJadi = penyelesaian
	details.PengambilanBarangBaku = pengambilan
	details.BarangMentah = barangMentah

	// Return the result as JSON
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    details})
}
