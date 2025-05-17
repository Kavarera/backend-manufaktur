package model

type BarangMentah struct {
	ID           int     `json:"id"`
	Nama         string  `json:"nama"`
	KodeBarang   string  `json:"kodeBarang"`
	HargaStandar float64 `json:"hargaStandar"`
	SatuanID     int     `json:"satuanId"`
	SatuanNama   string  `json:"satuanNama"` // joined from satuanTurunan.nama
	Stok         float64 `json:"stok"`
	GudangID     int     `json:"gudangId"`
	GudangNama   string  `json:"gudangNama"` // joined from gudang.nama
}
