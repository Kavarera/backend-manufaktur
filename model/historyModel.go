package model

type PerintahKerjaDetails struct {
	PerintahKerja          PerintahKerja            `json:"perintahKerja"`
	PenyelesaianBarangJadi []PenyelesaianBarangJadi `json:"penyelesaianBarangJadi"`
	PengambilanBarangBaku  []PengambilanBarangBaku  `json:"pengambilanBarangBaku"`
	BarangMentah           []BarangMentah           `json:"barangMentah"`
}
