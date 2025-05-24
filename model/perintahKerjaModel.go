package model

import "time"

type PerintahKerja struct {
	ID             string     `json:"id"`
	TanggalRilis   *time.Time `json:"tanggalRilis"`
	TanggalProgres *time.Time `json:"tanggalProgres"`
	TanggalSelesai *time.Time `json:"tanggalSelesai"`
	Status         string     `json:"status"`
	Hasil          float64    `json:"hasil"`
	Customer       *string    `json:"customer"`
	Keterangan     string     `json:"keterangan"`
}
