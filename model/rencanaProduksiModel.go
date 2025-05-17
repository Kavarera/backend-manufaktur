package model

import "time"

type RencanaProduksi struct {
	ID               *string    `json:"id,omitempty"`
	BarangProduksiID *int       `json:"barangProduksiId,omitempty"`
	TanggalMulai     *time.Time `json:"tanggalMulai,omitempty"`
	TanggalSelesai   *time.Time `json:"tanggalSelesai,omitempty"`
}

type RencanaProduksiAdd struct {
	ID               string `json:"id"` // or pointer if optional
	BarangProduksiID int    `json:"barangProduksiId"`

	TanggalMulai string `json:"tanggalMulai"` // e.g. "2025-06-01"
	WaktuMulai   string `json:"waktuMulai"`   // e.g. "08:00:00"

	TanggalSelesai string `json:"tanggalSelesai"` // e.g. "2025-06-15"
	WaktuSelesai   string `json:"waktuSelesai"`   // e.g. "17:00:00"
}
