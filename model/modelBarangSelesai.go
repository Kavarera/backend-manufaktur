package model

import (
	"fmt"
	"time"
)

// PenyelesaianBarangJadi represents the penyelesaian barang jadi record in the database
type PenyelesaianBarangJadi struct {
	ID                  int         `json:"id"`
	IDPerintahKerja     string      `json:"idPerintahKerja"`
	Nama                string      `json:"nama"`
	Jumlah              float64     `json:"jumlah"`
	TanggalPenyelesaian CustomDate2 `json:"tanggalPenyelesaian"`
}

const DateFormat2 = "2006-01-02" // Custom date format: yyyy-mm-dd

// CustomDate2 defines a type for parsing date-only values (yyyy-mm-dd).
type CustomDate2 time.Time

// UnmarshalJSON is used to parse date strings in yyyy-mm-dd format
func (d *CustomDate2) UnmarshalJSON(b []byte) error {
	// Removing the quotes from the date string (e.g., "2023-07-01" becomes 2023-07-01)
	str := string(b)
	if str == "\"\"" || str == "null" || str == "" {
		// Handle empty string as null value
		*d = CustomDate2{}
		return nil
	}
	str = str[1 : len(str)-1]                       // Remove quotes
	parsedTime, err := time.Parse(DateFormat2, str) // Fixed: use DateFormat2
	if err != nil {
		return fmt.Errorf("invalid date format: %s", err)
	}
	*d = CustomDate2(parsedTime)
	return nil
}

// MarshalJSON is used to output the date in a custom format (yyyy-mm-dd)
func (d *CustomDate2) MarshalJSON() ([]byte, error) {
	// If CustomDate is nil, return "null"
	if d == nil || time.Time(*d).IsZero() {
		return []byte("null"), nil
	}
	// Format the date into yyyy-mm-dd format and return it as a JSON string with quotes
	formatted := time.Time(*d).Format(DateFormat2) // Fixed: use DateFormat2
	return []byte(`"` + formatted + `"`), nil      // Fixed: add quotes
}

// ToTime converts CustomDate2 to time.Time, suitable for SQL operations
func (d *CustomDate2) ToTime() time.Time {
	if d == nil {
		return time.Time{}
	}
	return time.Time(*d)
}

// String returns string representation of the date
func (d *CustomDate2) String() string {
	if d == nil || time.Time(*d).IsZero() {
		return ""
	}
	return time.Time(*d).Format(DateFormat2)
}
