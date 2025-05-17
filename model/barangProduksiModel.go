package model

type BarangProduksi struct {
	ID           int     `json:"id"`
	Nama         string  `json:"nama"`
	KodeBarang   string  `json:"kodeBarang"`
	HargaStandar float64 `json:"hargaStandar"`
	HargaReal    float64 `json:"hargaReal"`
	SatuanID     int     `json:"satuanId"`
	SatuanNama   string  `json:"satuanNama"` // from satuanTurunan.nama
	Stok         float64 `json:"stok"`
	GudangID     int     `json:"gudangId"`
	GudangNama   string  `json:"gudangNama"` // from gudang.nama
}
