package model

type BarangMentah struct {
	ID              int     `json:"id"`
	Nama            string  `json:"nama"`
	KodeBarang      string  `json:"kodeBarang"`
	HargaStandar    float64 `json:"hargaStandar"`
	SatuanID        *int    `json:"satuanId"` // joined from satuanTurunan.nama
	SatuanNama      *string `json:"satuanNama"`
	SatuanUtamaID   *int    `json:"satuanUtamaId"`
	SatuanUtamaNama *string `json:"satuanUtamaNama"`
	Stok            float64 `json:"stok"`
	GudangID        int     `json:"gudangId"`
	GudangNama      string  `json:"gudangNama"` // joined from gudang.nama
}
