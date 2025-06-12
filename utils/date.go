package utils

import (
	"log"
	"time"
)

// ToTime converts a string in dd-mm-yyyy format to time.Time
func ToTime(dateStr string) (*time.Time, error) {
	// Parse the date in dd-mm-yyyy format
	layout := "02-01-2006" // dd-mm-yyyy
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Error parsing date:", err)
		return nil, err
	}
	return &t, nil
}

// FormatDate formats time.Time to dd-mm-yyyy format string
func FormatDate(t *time.Time) string {
	if t != nil {
		return t.Format("02-01-2006") // dd-mm-yyyy
	}
	return ""
}

func ListFormatDate(t time.Time) string {
	return t.Format("02-01-2006")
}
