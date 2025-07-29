package model

import "time"

type FormulaProduksi struct {
	ID               int       `json:"id"`
	IDBarangProduksi int       `json:"idBarangProduksi"`
	Kuantitas        float64   `json:"kuantitas"`
	TanggalMulai     time.Time `json:"tanggalMulai"`
	NamaProduksi     string    `json:"namaProduksi"`
	NamaFormula      string    `json:"namaFormula"`
}
