package handler

import (
	"database/sql"
	"fmt"
	"io"
	"manufacture_API/db"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadDocumentForPerintahKerja(c *gin.Context) {
	perintahKerjaID := c.Param("id")

	// Check if perintah kerja exists
	var exists bool
	err := db.GetDB().QueryRow(`SELECT EXISTS(SELECT 1 FROM "perintahKerja" WHERE "id" = $1)`, perintahKerjaID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Perintah Kerja not found",
		})
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Failed to get file from request",
		})
		return
	}
	defer file.Close()

	// Validate file type
	if filepath.Ext(strings.ToLower(header.Filename)) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Only PDF files are allowed",
		})
		return
	}

	// Validate file size (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "File size too large. Maximum 10MB allowed",
		})
		return
	}

	// Generate unique filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s_%s", perintahKerjaID, timestamp, header.Filename)

	// Create uploads directory if it doesn't exist
	uploadsDir := "./uploads/documents"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to create upload directory",
		})
		return
	}

	// Save file to local storage
	filePath := filepath.Join(uploadsDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to create file",
		})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to save file",
		})
		return
	}

	// Update database with document info
	documentURL := fmt.Sprintf("/uploads/documents/%s", filename)
	updateQuery := `UPDATE "perintahKerja" SET "document_url" = $1, "document_nama" = $2 WHERE "id" = $3`

	_, err = db.GetDB().Exec(updateQuery, documentURL, header.Filename, perintahKerjaID)
	if err != nil {
		// Clean up file if database update fails
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Failed to update document info in database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "OK",
		"message":     "Document uploaded successfully",
		"documentUrl": documentURL,
		"filename":    header.Filename,
	})
}

// DownloadDocument serves the document file
func DownloadDocument(c *gin.Context) {
	id := c.Param("id")

	var documentURL, documentName sql.NullString
	err := db.GetDB().QueryRow(
		`SELECT "document_url", "document_nama" FROM "perintahKerja" WHERE "id"=$1`,
		id,
	).Scan(&documentURL, &documentName)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Perintah Kerja not found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Database error",
		})
		return
	}

	if !documentURL.Valid || documentURL.String == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "No document found for this Perintah Kerja",
		})
		return
	}

	// Extract filename from URL
	parts := strings.Split(documentURL.String, "/")
	filename := parts[len(parts)-1]
	filePath := filepath.Join("./uploads/documents", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Error",
			"message": "Document file not found on server",
		})
		return
	}

	// Set headers for file download
	displayName := documentName.String
	if displayName == "" {
		displayName = filename
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))
	c.Header("Content-Type", "application/pdf")

	c.File(filePath)
}
