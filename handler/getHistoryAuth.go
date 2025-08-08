package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPerintahKerjaDetailsByID fetches Perintah Kerja along with related PenyelesaianBarangJadi, PengambilanBarangBaku, and BarangMentah data
func GetPerintahKerjaDetailsByID(c *gin.Context) {
	id := c.Param("id")

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
		pk.id_rencana_produksi,
		rp."namaProduksi",
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
		bm.stok AS barang_mentah_stok,
		bm.jumlah_turunan AS jumlah_satuan_turunan
	FROM 
		"perintahKerja" pk
	LEFT JOIN "rencanaProduksi" rp ON rp.id = pk.id_rencana_produksi
	LEFT JOIN "penyelesaianBarangJadi" pbd ON pbd.id_perintah_kerja = pk.id
	LEFT JOIN "pengambilanBarangBaku" pbb ON pbb.id_perintah_kerja = pk.id
	LEFT JOIN "barangMentah" bm ON bm.id = pbb.id_barang_mentah
	WHERE pk.id = $1
	`

	rows, err := db.GetDB().Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch perintah kerja details"})
		return
	}
	defer rows.Close()

	var (
		details      model.PerintahKerjaDetails
		penyelesaian []model.PenyelesaianBarangJadi
		pengambilan  []model.PengambilanBarangBaku
		barangMentah []model.BarangMentah
	)

	for rows.Next() {
		var pk model.PerintahKerja
		var pr model.RencanaProduksi
		var pbb model.PengambilanBarangBaku
		var bm model.BarangMentah
		var jm model.BarangSatuanTurunanMentah

		var (
			tanggalProgres      sql.NullString
			tanggalSelesai      sql.NullString
			jumlahTurunan       sql.NullFloat64
			penyelesaianID      sql.NullInt64
			namaPenyelesaian    sql.NullString
			jumlahPenyelesaian  sql.NullFloat64
			tanggalPenyelesaian sql.NullTime
		)

		err := rows.Scan(
			&pk.ID, &pk.TanggalRilis, &tanggalProgres, &tanggalSelesai, &pk.Status,
			&pk.Hasil, &pk.Customer, &pk.Keterangan, &pk.DocumentURL, &pk.DocumentNama, &pk.IdRencanaProduksi,
			&pr.NamaProduksi,
			&penyelesaianID, &namaPenyelesaian, &jumlahPenyelesaian, &tanggalPenyelesaian,
			&pbb.ID, &pbb.Kebutuhan,
			&bm.ID, &bm.Nama, &bm.KodeBarang, &bm.HargaStandar, &bm.Stok, &jumlahTurunan,
		)
		if err != nil {
			fmt.Println("Scan error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}

		pk.TanggalProgres = ""
		pk.TanggalSelesai = ""
		jm.JumlahTurunan = 0

		// Convert nullable time fields
		if tanggalProgres.Valid {
			pk.TanggalProgres = tanggalProgres.String
		}
		if tanggalSelesai.Valid {
			pk.TanggalSelesai = tanggalSelesai.String
		}
		if jumlahTurunan.Valid {
			jm.JumlahTurunan = jumlahTurunan.Float64
		}

		// Fill the penyelesaian data
		if penyelesaianID.Valid {
			pbd := model.PenyelesaianBarangJadi{
				ID: int(penyelesaianID.Int64),
			}
			if namaPenyelesaian.Valid {
				pbd.Nama = namaPenyelesaian.String
			}
			if jumlahPenyelesaian.Valid {
				pbd.Jumlah = jumlahPenyelesaian.Float64
			}
			if tanggalPenyelesaian.Valid {
				date := model.CustomDate2(tanggalPenyelesaian.Time)
				pbd.TanggalPenyelesaian = date
			}
			penyelesaian = append(penyelesaian, pbd)
		}

		if pbb.ID != 0 {
			pengambilan = append(pengambilan, pbb)
		}
		if bm.ID != 0 {
			barangMentah = append(barangMentah, bm)
		}

		details.PerintahKerja = pk
	}

	details.PenyelesaianBarangJadi = penyelesaian
	details.PengambilanBarangBaku = pengambilan
	details.BarangMentah = barangMentah

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    details,
	})
}

func GetPerintahKerjaDetails(c *gin.Context) {
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
		pk.id_rencana_produksi,
		rp."namaProduksi",
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
		bm.stok AS barang_mentah_stok,
		bm.jumlah_turunan AS jumlah_satuan_turunan
	FROM 
		"perintahKerja" pk
	LEFT JOIN "rencanaProduksi" rp ON rp.id = pk.id_rencana_produksi
	LEFT JOIN "penyelesaianBarangJadi" pbd ON pbd.id_perintah_kerja = pk.id
	LEFT JOIN "pengambilanBarangBaku" pbb ON pbb.id_perintah_kerja = pk.id
	LEFT JOIN "barangMentah" bm ON bm.id = pbb.id_barang_mentah
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch perintah kerja details"})
		return
	}
	defer rows.Close()

	var (
		detailsMap = map[string]*model.PerintahKerjaDetails{}
		result     []model.PerintahKerjaDetails
	)

	for rows.Next() {
		var (
			pk                   model.PerintahKerja
			pr                   model.RencanaProduksi
			pbb                  model.PengambilanBarangBaku
			bm                   model.BarangMentah
			tanggalProgres       sql.NullString
			tanggalSelesai       sql.NullString
			jumlahTurunan        sql.NullFloat64
			penyelesaianID       sql.NullInt64
			namaPenyelesaian     sql.NullString
			jumlahPenyelesaian   sql.NullFloat64
			tanggalPenyelesaian  sql.NullTime
			pengambilanID        sql.NullInt64
			pengambilanKebutuhan sql.NullFloat64
			mentahID             sql.NullInt64
			mentahNama           sql.NullString
			mentahKode           sql.NullString
			mentahHarga          sql.NullFloat64
			mentahStok           sql.NullFloat64
		)

		err := rows.Scan(
			&pk.ID, &pk.TanggalRilis, &tanggalProgres, &tanggalSelesai, &pk.Status,
			&pk.Hasil, &pk.Customer, &pk.Keterangan, &pk.DocumentURL, &pk.DocumentNama, &pk.IdRencanaProduksi,
			&pr.NamaProduksi,
			&penyelesaianID, &namaPenyelesaian, &jumlahPenyelesaian, &tanggalPenyelesaian,
			&pengambilanID, &pengambilanKebutuhan,
			&mentahID, &mentahNama, &mentahKode, &mentahHarga, &mentahStok, &jumlahTurunan,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}

		if tanggalProgres.Valid {
			pk.TanggalProgres = tanggalProgres.String
		}
		if tanggalSelesai.Valid {
			pk.TanggalSelesai = tanggalSelesai.String
		}

		// Group by PerintahKerja.ID
		group, exists := detailsMap[pk.ID]
		if !exists {
			group = &model.PerintahKerjaDetails{
				PerintahKerja:          pk,
				PenyelesaianBarangJadi: []model.PenyelesaianBarangJadi{},
				PengambilanBarangBaku:  []model.PengambilanBarangBaku{},
				BarangMentah:           []model.BarangMentah{},
			}
			detailsMap[pk.ID] = group
		}

		// Penyelesaian
		if penyelesaianID.Valid {
			pbd := model.PenyelesaianBarangJadi{ID: int(penyelesaianID.Int64)}
			if namaPenyelesaian.Valid {
				pbd.Nama = namaPenyelesaian.String
			}
			if jumlahPenyelesaian.Valid {
				pbd.Jumlah = jumlahPenyelesaian.Float64
			}
			if tanggalPenyelesaian.Valid {
				date := model.CustomDate2(tanggalPenyelesaian.Time)
				pbd.TanggalPenyelesaian = date
			}
			group.PenyelesaianBarangJadi = append(group.PenyelesaianBarangJadi, pbd)
		}

		// Pengambilan
		if pengambilanID.Valid {
			pbb.ID = int(pengambilanID.Int64)
			if pengambilanKebutuhan.Valid {
				pbb.Kebutuhan = pengambilanKebutuhan.Float64
			}
			group.PengambilanBarangBaku = append(group.PengambilanBarangBaku, pbb)
		}

		// Barang Mentah
		if mentahID.Valid {
			bm.ID = int(mentahID.Int64)
			if mentahNama.Valid {
				bm.Nama = mentahNama.String
			}
			if mentahKode.Valid {
				bm.KodeBarang = mentahKode.String
			}
			if mentahHarga.Valid {
				bm.HargaStandar = mentahHarga.Float64
			}
			if mentahStok.Valid {
				bm.Stok = mentahStok.Float64
			}
			if jumlahTurunan.Valid {
				bm.SatuanTurunan = []model.BarangSatuanTurunanMentah{
					{JumlahTurunan: jumlahTurunan.Float64},
				}
			}
			group.BarangMentah = append(group.BarangMentah, bm)
		}
	}

	// Collect results
	for _, v := range detailsMap {
		result = append(result, *v)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    result,
	})
}
