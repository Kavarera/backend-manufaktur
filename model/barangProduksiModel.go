package model

type BarangProduksi struct {
	ID              int                   `json:"id"`
	Nama            string                `json:"nama"`
	KodeBarang      string                `json:"kodeBarang"`
	HargaStandar    float64               `json:"hargaStandar"`
	HargaReal       float64               `json:"hargaReal"`
	SatuanTurunan   []BarangSatuanTurunan `json:"satuanTurunan"`
	Stok            float64               `json:"stok"`
	GudangID        int                   `json:"gudangId"`
	GudangNama      string                `json:"gudangNama"`
	SatuanUtamaID   *int                  `json:"satuanUtamaId"`
	SatuanUtamaNama *string               `json:"satuanUtamaNama"`
}

type BarangSatuanTurunan struct {
	SatuanTurunanID   int     `json:"satuanTurunanId"`
	SatuanTurunanNama string  `json:"satuanTurunanNama"`
	Jumlah            float64 `json:"jumlah"`
}
