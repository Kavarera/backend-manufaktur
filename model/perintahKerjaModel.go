package model

import (
	"database/sql/driver"
	"time"
)

type PerintahKerja struct {
	ID             string     `json:"id"`
	TanggalRilis   *time.Time `json:"tanggalRilis"`
	TanggalProgres *time.Time `json:"tanggalProgres"`
	TanggalSelesai *time.Time `json:"tanggalSelesai"`
	Status         string     `json:"status"`
	Hasil          float64    `json:"hasil"`
	Customer       *string    `json:"customer"`
	Keterangan     string     `json:"keterangan"`
	DocumentURL    *string    `json:"documentUrl,omitempty"`
	DocumentNama   *string    `json:"documentNama,omitempty"`
}

// For file upload request
type PerintahKerjaWithFile struct {
	PerintahKerja
	HasDocument bool `json:"has_document"`
}

// Custom scanner for nullable time
type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	nt.Time, nt.Valid = value.(time.Time), true
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
