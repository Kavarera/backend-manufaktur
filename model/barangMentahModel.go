package model

type BarangSatuanTurunanMentah struct {
	SatuanID      int     `json:"satuanId"`
	SatuanNama    string  `json:"satuanNama"`
	JumlahTurunan float64 `json:"jumlahTurunan"`
}

type BarangMentah struct {
	ID              int                         `json:"id"`
	Nama            string                      `json:"nama"`
	KodeBarang      string                      `json:"kodeBarang"`
	HargaStandar    float64                     `json:"hargaStandar"`
	SatuanID        *int                        `json:"satuanId"`
	SatuanUtamaID   *int                        `json:"satuanUtamaId"`
	SatuanUtamaNama *string                     `json:"satuanUtamaNama"`
	Stok            float64                     `json:"stok"`
	GudangID        int                         `json:"gudangId"`
	GudangNama      string                      `json:"gudangNama"`
	SatuanTurunan   []BarangSatuanTurunanMentah `json:"satuanTurunan"`
}
