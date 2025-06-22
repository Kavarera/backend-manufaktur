package model

import "time"

type PengambilanBarangBaku struct {
	ID                       int        `json:"id"`
	IDPerintahKerja          string     `json:"idPerintahKerja"`
	IDBarangMentah           int        `json:"idBarangMentah"`
	Kebutuhan                float64    `json:"kebutuhan"`
	TanggalRilis             time.Time  `json:"tanggalRilis"`
	TanggalProgres           *time.Time `json:"tanggalProgres"`
	TanggalSelesai           *time.Time `json:"tanggalSelesai"`
	StatusPerintahKerja      string     `json:"statusPerintahKerja"`
	NamaBarangMentah         string     `json:"namaBarangMentah"`
	KodeBarangMentah         string     `json:"kodeBarangMentah"`
	HargaStandarBarangMentah float64    `json:"hargaStandarBarangMentah"`
	StokBarangMentah         float64    `json:"stokBarangMentah"`
	TanggalWaktu             time.Time  `json:"tanggalWaktu"`
}
