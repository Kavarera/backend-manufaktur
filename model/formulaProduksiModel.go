package model

// FormulaProduksi represents the production formula in the database
type FormulaProduksi struct {
	ID            int     `json:"id"`
	BarangJadi    string  `json:"barang_jadi"`
	Kuantitas     float64 `json:"kuantitas"`
	Satuan        float64 `json:"satuan"`
	BahanBaku     string  `json:"bahan_baku"`
	SatuanTurunan float64 `json:"satuanTurunan"`
}
