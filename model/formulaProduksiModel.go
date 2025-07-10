package model

type FormulaProduksi struct {
	ID            int     `json:"id"`
	BarangJadi    string  `json:"barangJadi"`
	Kuantitas     float64 `json:"kuantitas"`
	Satuan        float64 `json:"satuan"`
	BahanBaku     string  `json:"bahanBaku"`
	SatuanTurunan float64 `json:"satuanTurunan"`
}
