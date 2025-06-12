package model

import (
	"fmt"
	"time"
)

type RencanaProduksiUpdate struct {
	ID               string      `json:"id"`
	BarangProduksiID *int        `json:"barangProduksiId,omitempty"`
	TanggalMulai     *CustomDate `json:"tanggalMulai,omitempty"`
	TanggalSelesai   *CustomDate `json:"tanggalSelesai,omitempty"`
}

// RencanaProduksi represents a production plan record in the database.
type RencanaProduksi struct {
	ID               string     `json:"id"`
	BarangProduksiID int        `json:"barangProduksiId"`
	TanggalMulai     CustomDate `json:"tanggalMulai"`
	TanggalSelesai   CustomDate `json:"tanggalSelesai,omitempty"`
}

const DateFormat = "2006-01-02"

// CustomDate defines a type for parsing date-only values (yyyy-mm-dd).
type CustomDate time.Time

// UnmarshalJSON is used to parse date strings in yyyy-mm-dd format
func (d *CustomDate) UnmarshalJSON(b []byte) error {
	// Remove quotes from the date string (e.g. "2023-07-01" becomes 2023-07-01)
	str := string(b)
	if len(str) < 2 {
		return fmt.Errorf("invalid date string: too short")
	}

	// Handle null values
	if str == "null" {
		return nil
	}

	// Remove quotes
	str = str[1 : len(str)-1]

	parsedTime, err := time.Parse(DateFormat, str)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", err)
	}

	*d = CustomDate(parsedTime)
	return nil
}

// MarshalJSON is used to output the date in a custom format (yyyy-mm-dd)
func (d *CustomDate) MarshalJSON() ([]byte, error) {
	// If CustomDate is nil or zero time, return "null"
	t := time.Time(*d)
	if t.IsZero() {
		return []byte("null"), nil
	}

	// Format the date into yyyy-mm-dd format and return it as a JSON string with quotes
	return []byte(`"` + t.Format(DateFormat) + `"`), nil
}

// ToTime converts CustomDate to time.Time, suitable for SQL operations
func (d *CustomDate) ToTime() time.Time {
	if d == nil {
		return time.Time{}
	}
	return time.Time(*d)
}

// Scan implements the sql.Scanner interface for database scanning
func (d *CustomDate) Scan(value interface{}) error {
	if value == nil {
		*d = CustomDate(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = CustomDate(v)
		return nil
	case string:
		t, err := time.Parse(DateFormat, v)
		if err != nil {
			return err
		}
		*d = CustomDate(t)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into CustomDate", value)
	}
}
