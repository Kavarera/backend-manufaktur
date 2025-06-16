package model

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type PerintahKerja struct {
	ID                      string       `json:"id"`
	TanggalRilis            string       `json:"tanggalRilis"`
	TanggalProgres          string       `json:"tanggalProgres"`
	TanggalSelesai          string       `json:"tanggalSelesai"`
	Status                  string       `json:"status"`
	Hasil                   float64      `json:"hasil"`
	Customer                *string      `json:"customer"`
	Keterangan              string       `json:"keterangan"`
	DocumentURL             *string      `json:"documentUrl,omitempty"`
	DocumentNama            *string      `json:"documentNama,omitempty"`
	TanggalRilisFormatted   string       `json:"tanggalRilisFormatted,omitempty"`
	TanggalProgresFormatted string       `json:"tanggalProgresFormatted,omitempty"`
	TanggalSelesaiFormatted string       `json:"tanggalSelesaiFormatted,omitempty"`
	IdRencanaProduksi       string       `json:"idRencanaProduksi"`
	NamaProduksi            string       `json:"namaProduksi"`
	TanggalRilisTime        *time.Time   `json:"-"` // This will not be part of the response JSON
	TanggalProgresTime      *time.Time   `json:"-"` // This will not be part of the response JSON
	TanggalSelesaiTime      *time.Time   `json:"-"`
	TanggalRilisTime2       sql.NullTime `json:"-"` // This will not be part of the response JSON
	TanggalProgresTime2     sql.NullTime `json:"-"` // This will not be part of the response JSON
	TanggalSelesaiTime2     sql.NullTime `json:"-"`
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
