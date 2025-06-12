package utils

import "strings"

var AllowedStatuses = []string{"Dijadwalkan", "Dalam Proses", "Selesai"}

// IsValidStatus checks if a given status is allowed
func IsValidStatus(status string) bool {
	status = strings.TrimSpace(strings.ToLower(status))
	for _, s := range AllowedStatuses {
		if strings.ToLower(s) == status {
			return true
		}
	}
	return false
}
