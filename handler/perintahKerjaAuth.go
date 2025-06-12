package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"manufacture_API/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListPerintahKerja returns all records
func ListPerintahKerja(c *gin.Context) {
	query := `
		SELECT "id", "tanggal_rilis", "tanggal_progres", "tanggal_selesai", 
		       "status", "hasil", "customer", "keterangan", "document_url", "document_nama"
		FROM "perintahKerja" 
		ORDER BY "id"
	`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to fetch data: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	var list []model.PerintahKerja
	for rows.Next() {
		var data model.PerintahKerja
		err := rows.Scan(
			&data.ID, &data.TanggalRilisTime2, &data.TanggalProgresTime2, &data.TanggalSelesaiTime2,
			&data.Status, &data.Hasil, &data.Customer, &data.Keterangan,
			&data.DocumentURL, &data.DocumentNama,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Error",
				"message": "Failed to parse data: " + err.Error(),
			})
			return
		}

		// Convert sql.NullTime to *time.Time (if valid)
		if data.TanggalRilisTime2.Valid {
			data.TanggalRilis = utils.ListFormatDate(data.TanggalRilisTime2.Time)
		}
		if data.TanggalProgresTime2.Valid {
			data.TanggalProgres = utils.ListFormatDate(data.TanggalProgresTime2.Time)
		}
		if data.TanggalSelesaiTime2.Valid {
			data.TanggalSelesai = utils.ListFormatDate(data.TanggalSelesaiTime2.Time)
		}

		list = append(list, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Berhasil",
		"data":    list,
	})
}

// AddPerintahKerja adds a new record
func AddPerintahKerja(c *gin.Context) {
	var input model.PerintahKerja
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	// Convert string dates to *time.Time
	var err error
	if input.TanggalRilis != "" {
		input.TanggalRilisTime, err = utils.ToTime(input.TanggalRilis) // Convert string to *time.Time
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid tanggal_rilis format"})
			return
		}
	}
	if input.TanggalProgres != "" {
		input.TanggalProgresTime, err = utils.ToTime(input.TanggalProgres) // Convert string to *time.Time
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid tanggal_progres format"})
			return
		}
	}
	if input.TanggalSelesai != "" {
		input.TanggalSelesaiTime, err = utils.ToTime(input.TanggalSelesai) // Convert string to *time.Time
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid tanggal_selesai format"})
			return
		}
	}
	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	// Validate status
	if !utils.IsValidStatus(input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Status must be one of: Dijadwalkan, Dalam Proses, Selesai",
		})
		return
	}

	// Proceed with database insertion
	query := `
		INSERT INTO "perintahKerja" 
		("id", "tanggal_rilis", "tanggal_progres", "tanggal_selesai", "status", "hasil", "customer", "keterangan", "document_url", "document_nama")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING "id"
	`

	var returnedID string
	err = db.GetDB().QueryRow(query,
		input.ID, input.TanggalRilisTime, input.TanggalProgresTime, input.TanggalSelesaiTime,
		input.Status, input.Hasil, input.Customer, input.Keterangan,
		input.DocumentURL, input.DocumentNama,
	).Scan(&returnedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to insert data: " + err.Error(),
		})
		return
	}

	input.ID = returnedID
	c.JSON(http.StatusCreated, gin.H{
		"status":  "OK",
		"message": "Berhasil menambah perintah kerja",
		"data":    input,
	})
}

// UpdatePerintahKerja updates by ID
func UpdatePerintahKerja(c *gin.Context) {
	id := c.Param("id")

	// Step 1: Retrieve the current record from the database
	var existingRecord model.PerintahKerja
	err := db.GetDB().QueryRow(`
		SELECT "id", "tanggal_rilis", "tanggal_progres", "tanggal_selesai", 
		       "status", "hasil", "customer", "keterangan", "document_url", "document_nama"
		FROM "perintahKerja"
		WHERE "id" = $1
	`, id).Scan(
		&existingRecord.ID, &existingRecord.TanggalRilisTime, &existingRecord.TanggalProgresTime,
		&existingRecord.TanggalSelesaiTime, &existingRecord.Status, &existingRecord.Hasil,
		&existingRecord.Customer, &existingRecord.Keterangan, &existingRecord.DocumentURL,
		&existingRecord.DocumentNama,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": "Perintah kerja not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "Failed to fetch existing record: " + err.Error()})
		return
	}

	// Step 2: Bind the incoming JSON to the handler's struct
	var pk model.PerintahKerja
	if err := c.ShouldBindJSON(&pk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Step 3: Handle the date fields (convert them only if provided)
	var updateValues []interface{}
	var setClauses []string
	argPos := 1

	// Handle tanggal_rilis (update only if new value is provided)
	if pk.TanggalRilis != "" {
		// Convert the string (dd-mm-yyyy) to *time.Time
		var err error
		pk.TanggalRilisTime, err = utils.ToTime(pk.TanggalRilis)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid tanggal_rilis format"})
			return
		}
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_rilis"=$%d`, argPos))
		updateValues = append(updateValues, pk.TanggalRilisTime)
		argPos++
	} else {
		// If no new date is provided, retain the existing one
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_rilis"=$%d`, argPos))
		updateValues = append(updateValues, existingRecord.TanggalRilisTime)
		argPos++
	}

	// Handle tanggal_progres (update only if new value is provided)
	if pk.TanggalProgres != "" {
		var err error
		pk.TanggalProgresTime, err = utils.ToTime(pk.TanggalProgres)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid tanggal_progres format"})
			return
		}
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_progres"=$%d`, argPos))
		updateValues = append(updateValues, pk.TanggalProgresTime)
		argPos++
	} else {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_progres"=$%d`, argPos))
		updateValues = append(updateValues, existingRecord.TanggalProgresTime)
		argPos++
	}

	// Handle tanggal_selesai (update only if new value is provided)
	if pk.TanggalSelesai != "" {
		var err error
		pk.TanggalSelesaiTime, err = utils.ToTime(pk.TanggalSelesai)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Invalid tanggal_selesai format"})
			return
		}
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_selesai"=$%d`, argPos))
		updateValues = append(updateValues, pk.TanggalSelesaiTime)
		argPos++
	} else {
		setClauses = append(setClauses, fmt.Sprintf(`"tanggal_selesai"=$%d`, argPos))
		updateValues = append(updateValues, existingRecord.TanggalSelesaiTime)
		argPos++
	}

	// Handle other fields: status, hasil, customer, etc.
	if pk.Status != "" {
		if !utils.IsValidStatus(pk.Status) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "Error",
				"message": "Status must be one of: Dijadwalkan, Dalam Proses, atau Selesai",
			})
			return
		}
		setClauses = append(setClauses, fmt.Sprintf(`"status"=$%d`, argPos))
		updateValues = append(updateValues, pk.Status)
		argPos++
	}
	if pk.Hasil != 0 {
		setClauses = append(setClauses, fmt.Sprintf(`"hasil"=$%d`, argPos))
		updateValues = append(updateValues, pk.Hasil)
		argPos++
	}
	if pk.Customer != nil {
		setClauses = append(setClauses, fmt.Sprintf(`"customer"=$%d`, argPos))
		updateValues = append(updateValues, pk.Customer)
		argPos++
	}
	if pk.Keterangan != "" {
		setClauses = append(setClauses, fmt.Sprintf(`"keterangan"=$%d`, argPos))
		updateValues = append(updateValues, pk.Keterangan)
		argPos++
	}

	// If no fields are provided for update
	if len(setClauses) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "No fields to update",
		})
		return
	}

	// Step 4: Update the record in the database
	updateValues = append(updateValues, id) // Add the ID to the end of the values for the WHERE clause
	query := fmt.Sprintf(`UPDATE "perintahKerja" SET %s WHERE "id"=$%d`, strings.Join(setClauses, ", "), argPos)

	res, err := db.GetDB().Exec(query, updateValues...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to update perintah kerja: " + err.Error(),
		})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Perintah kerja not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Perintah kerja updated successfully",
	})
}

// DeletePerintahKerja deletes by ID
func DeletePerintahKerja(c *gin.Context) {
	id := c.Param("id")

	// First, get document info if exists
	var documentURL sql.NullString
	err := db.GetDB().QueryRow(`SELECT "document_url" FROM "perintahKerja" WHERE "id"=$1`, id).Scan(&documentURL)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to check document info",
		})
		return
	}

	// Delete from database
	res, err := db.GetDB().Exec(`DELETE FROM "perintahKerja" WHERE "id"=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to delete data: " + err.Error(),
		})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Data not found",
		})
		return
	}

	// Delete associated file if exists
	if documentURL.Valid && documentURL.String != "" {
		// Extract filename from URL
		parts := strings.Split(documentURL.String, "/")
		if len(parts) > 0 {
			filename := parts[len(parts)-1]
			filePath := filepath.Join("./uploads/documents", filename)
			os.Remove(filePath) // Ignore error if file doesn't exist
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Data berhasil dihapus",
	})
}
