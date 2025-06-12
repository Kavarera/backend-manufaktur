package model

type PengambilanBarangBaku struct {
	ID              int     `json:"id"`
	IDPerintahKerja string  `json:"idPerintahKerja"`
	IDBarangMentah  int     `json:"idBarangMentah"`
	Kebutuhan       float64 `json:"kebutuhan"`
}
