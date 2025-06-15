package handler

import (
	"database/sql"
	"fmt"
	"manufacture_API/db"
	"manufacture_API/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Standardized error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Standardized success response
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Status:  "Error",
		Message: message,
	})
}

func sendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Status:  "OK",
		Message: message,
	}
	if data != nil {
		response.Data = data
	}
	c.JSON(statusCode, response)
}

func ListRencanaProduksi(c *gin.Context) {
	query := `SELECT "id", "id_barang_produksi", "tanggal_mulai", "tanggal_selesai", "namaProduksi", "quantity" FROM "rencanaProduksi" ORDER BY "id"`

	rows, err := db.GetDB().Query(query)
	if err != nil {
		fmt.Printf("Error fetching rencana produksi: %v\n", err)
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rencana produksi")
		return
	}
	defer rows.Close()

	var list []model.RencanaProduksi
	for rows.Next() {
		var rp model.RencanaProduksi
		err := rows.Scan(&rp.ID, &rp.BarangProduksiID, &rp.TanggalMulai, &rp.TanggalSelesai, &rp.NamaProduksi, &rp.Quantity)
		if err != nil {
			fmt.Printf("Error scanning rencana produksi: %v\n", err)
			sendErrorResponse(c, http.StatusInternalServerError, "Failed to parse rencana produksi")
			return
		}
		list = append(list, rp)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		fmt.Printf("Error during rows iteration: %v\n", err)
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rencana produksi")
		return
	}

	sendSuccessResponse(c, http.StatusOK, "Berhasil", list)
}

func GetRencanaProduksiByID(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	query := `SELECT "id", "id_barang_produksi", "tanggal_mulai", "tanggal_selesai", "namaProduksi", "quantity" FROM "rencanaProduksi" WHERE "id"=$1`
	row := db.GetDB().QueryRow(query, id)

	var rp model.RencanaProduksi
	err := row.Scan(&rp.ID, &rp.BarangProduksiID, &rp.TanggalMulai, &rp.TanggalSelesai, &rp.NamaProduksi, &rp.Quantity)
	if err == sql.ErrNoRows {
		sendErrorResponse(c, http.StatusNotFound, "Rencana produksi not found")
		return
	} else if err != nil {
		fmt.Printf("Error fetching rencana produksi: %v\n", err)
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rencana produksi")
		return
	}

	sendSuccessResponse(c, http.StatusOK, "Berhasil", rp)
}

func AddRencanaProduksi(c *gin.Context) {
	var input model.RencanaProduksi
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Generate UUID for the id
	input.ID = uuid.New().String()

	// Convert CustomDate to time.Time before inserting into database
	tanggalMulai := input.TanggalMulai.ToTime()
	tanggalSelesai := input.TanggalSelesai.ToTime()

	// Validate dates (optional business logic)
	if tanggalSelesai.Before(tanggalMulai) {
		sendErrorResponse(c, http.StatusBadRequest, "Tanggal selesai cannot be before tanggal mulai")
		return
	}

	// Insert rencanaProduksi
	query := `
		INSERT INTO "rencanaProduksi" ("id", "id_barang_produksi", "tanggal_mulai", "tanggal_selesai", "namaProduksi", "quantity")
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.GetDB().Exec(query, input.ID, input.BarangProduksiID, tanggalMulai, tanggalSelesai, input.NamaProduksi, input.Quantity)
	if err != nil {
		fmt.Printf("Error inserting rencana produksi: %v\n", err)
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to add rencana produksi")
		return
	}

	sendSuccessResponse(c, http.StatusCreated, "Rencana Produksi created successfully", input)
}

func UpdateRencanaProduksi(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var payload model.RencanaProduksiUpdate
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updates := make(map[string]interface{})

	if payload.BarangProduksiID != nil {
		updates["id_barang_produksi"] = *payload.BarangProduksiID
	}

	if payload.TanggalMulai != nil {
		updates["tanggal_mulai"] = payload.TanggalMulai.ToTime()
	}

	if payload.TanggalSelesai != nil {
		updates["tanggal_selesai"] = payload.TanggalSelesai.ToTime()
	}

	if payload.NamaProduksi != nil {
		updates["namaProduksi"] = *payload.NamaProduksi
	}

	if payload.Quantity != nil {
		updates["quantity"] = *payload.Quantity
	}

	if len(updates) == 0 {
		sendErrorResponse(c, http.StatusBadRequest, "No fields to update")
		return
	}

	if payload.TanggalMulai != nil && payload.TanggalSelesai != nil {
		if payload.TanggalSelesai.ToTime().Before(payload.TanggalMulai.ToTime()) {
			sendErrorResponse(c, http.StatusBadRequest, "Tanggal selesai cannot be before tanggal mulai")
			return
		}
	}

	// Build dynamic SQL
	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argPos := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf(`"%s"=$%d`, column, argPos))
		args = append(args, value)
		argPos++
	}

	args = append(args, id)
	sql := fmt.Sprintf(`UPDATE "rencanaProduksi" SET %s WHERE "id"=$%d`, strings.Join(setClauses, ", "), argPos)

	res, err := db.GetDB().Exec(sql, args...)
	if err != nil {
		fmt.Printf("Error updating rencana produksi: %v\n", err)
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to update rencana produksi")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		sendErrorResponse(c, http.StatusNotFound, "Rencana produksi not found")
		return
	}

	sendSuccessResponse(c, http.StatusOK, "Rencana produksi updated successfully", nil)
}

func DeleteRencanaProduksi(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	query := `DELETE FROM "rencanaProduksi" WHERE "id"=$1`

	res, err := db.GetDB().Exec(query, id)
	if err != nil {
		fmt.Printf("Error deleting rencana produksi: %v\n", err)
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to delete rencana produksi")
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		sendErrorResponse(c, http.StatusNotFound, "Rencana produksi not found")
		return
	}

	sendSuccessResponse(c, http.StatusOK, "Rencana produksi deleted successfully", nil)
}
