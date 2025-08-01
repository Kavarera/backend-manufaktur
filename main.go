package main

import (
	"log"
	"manufacture_API/db"
	"manufacture_API/handler"
	"manufacture_API/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB connection
	db.InitDB()

	// Create Gin router with default middleware
	r := gin.Default()

	// Public routes (no authentication)
	r.POST("/login", handler.Login)

	// User Management Routes
	r.POST("/register", middleware.PermissionMiddleware("users:create"), handler.Register)
	r.GET("/users", middleware.PermissionMiddleware("users:read"), handler.AllUserList)
	r.GET("/users/:username", middleware.PermissionMiddleware("users:read"), handler.UserList)
	r.GET("/users/:username/roles", middleware.PermissionMiddleware("users:read"), handler.GetUserRoles)
	r.PUT("/users/:username/roles", middleware.PermissionMiddleware("users:update"), handler.UpdateUserRoles)
	r.DELETE("/users/:username", middleware.PermissionMiddleware("users:delete"), handler.UserDelete)

	r.GET("/barangSelesai", middleware.PermissionMiddleware("selesai:read"), handler.GetPenyelesaianBarangJadi)
	r.POST("/barangSelesai", middleware.PermissionMiddleware("selesai:create"), handler.AddPenyelesaianBarangJadi)
	r.PUT("/barangSelesai/:id", middleware.PermissionMiddleware("selesai:update"), handler.UpdatePenyelesaianBarangJadi)
	r.DELETE("/barangSelesai/:id", middleware.PermissionMiddleware("selesai:delete"), handler.DeletePenyelesaianBarangJadi)

	// Barang Produksi Routes
	r.GET("/barangProduksi", middleware.PermissionMiddleware("barang:read"), handler.ListBarangProduksi)
	r.GET("/barangProduksi/:id", middleware.PermissionMiddleware("barang:read"), handler.GetBarangProduksiByID)
	r.POST("/barangProduksi", middleware.PermissionMiddleware("barang:create"), handler.AddBarangProduksi)
	r.PUT("/barangProduksi/:id", middleware.PermissionMiddleware("barang:update"), handler.UpdateBarangProduksi)
	r.DELETE("/barangProduksi/:id", middleware.PermissionMiddleware("barang:delete"), handler.DeleteBarangProduksi)

	// Gudang Routes
	r.GET("/gudang", middleware.PermissionMiddleware("gudang:read"), handler.ListGudang)
	r.GET("/gudang/:id", middleware.PermissionMiddleware("gudang:read"), handler.GetGudangByID)
	r.POST("/gudang", middleware.PermissionMiddleware("gudang:create"), handler.AddGudang)
	r.PUT("/gudang/:id", middleware.PermissionMiddleware("gudang:update"), handler.UpdateGudang)
	r.DELETE("/gudang/:id", middleware.PermissionMiddleware("gudang:delete"), handler.DeleteGudang)

	// Barang Mentah Routes
	r.GET("/barangMentah", middleware.PermissionMiddleware("mentah:read"), handler.ListMentah)
	r.POST("/barangMentah", middleware.PermissionMiddleware("mentah:create"), handler.AddMentah)
	r.PUT("/barangMentah/:id", middleware.PermissionMiddleware("mentah:update"), handler.UpdateMentah)
	r.DELETE("/barangMentah/:id", middleware.PermissionMiddleware("mentah:delete"), handler.DeleteMentah)

	// Rencana Produksi Routes
	r.GET("/rencanaProduksi", middleware.PermissionMiddleware("rencana:read"), handler.ListRencanaProduksi)
	r.GET("/rencanaProduksi/:id", middleware.PermissionMiddleware("rencana:read"), handler.GetRencanaProduksiByID)
	r.POST("/rencanaProduksi", middleware.PermissionMiddleware("rencana:create"), handler.AddRencanaProduksi)
	r.PUT("/rencanaProduksi/:id", middleware.PermissionMiddleware("rencana:update"), handler.UpdateRencanaProduksi)
	r.DELETE("/rencanaProduksi/:id", middleware.PermissionMiddleware("rencana:delete"), handler.DeleteRencanaProduksi)
	r.GET("/jadwalProduksi", middleware.PermissionMiddleware("jadwal:read"), handler.ListRencanaProduksi)

	// Perintah Kerja Routes
	r.GET("/perintahKerja", middleware.PermissionMiddleware("perintah:read"), handler.ListPerintahKerja)
	r.POST("/perintahKerja", middleware.PermissionMiddleware("perintah:create"), handler.AddPerintahKerja)
	r.PUT("/perintahKerja/:id", middleware.PermissionMiddleware("perintah:update"), handler.UpdatePerintahKerja)
	r.DELETE("/perintahKerja/:id", middleware.PermissionMiddleware("perintah:delete"), handler.DeletePerintahKerja)
	r.POST("/perintahKerja/:id/upload-document", middleware.PermissionMiddleware("perintah:update"), handler.UploadDocumentForPerintahKerja)
	r.GET("/perintahKerja/:id/download-document", middleware.PermissionMiddleware("perintah:read"), handler.DownloadDocument)
	r.PUT("/updatePengerjaan/:id", middleware.PermissionMiddleware("perintah:update"), handler.UpdateProsesPengerjaan)

	// Pengambilan Barang Baku Routes
	r.GET("/pengambilanBarangBaku", middleware.PermissionMiddleware("pengambilan:read"), handler.GetPengambilanBarangBaku)
	r.POST("/pengambilanBarangBaku", middleware.PermissionMiddleware("pengambilan:create"), handler.AddPengambilanBarangBaku)
	r.PUT("/pengambilanBarangBaku/:id", middleware.PermissionMiddleware("pengambilan:update"), handler.UpdatePengambilanBarangBaku)
	r.DELETE("/pengambilanBarangBaku/:id", middleware.PermissionMiddleware("pengambilan:delete"), handler.DeletePengambilanBarangBaku)

	//History Routes
	r.GET("/history/:id", middleware.PermissionMiddleware("history:read"), handler.GetPerintahKerjaDetailsByID)

	//Formula Produksi Routes
	r.GET("/formulaProduksi", middleware.PermissionMiddleware("formula:read"), handler.ListFormulaProduksi)
	r.GET("/formulaProduksi/:id", middleware.PermissionMiddleware("formula:read"), handler.GetFormulaProduksiByID)
	r.POST("/formulaProduksi", middleware.PermissionMiddleware("formula:create"), handler.AddFormulaProduksi)
	r.PUT("/formulaProduksi/:id", middleware.PermissionMiddleware("formula:update"), handler.UpdateFormulaProduksi)
	r.DELETE("/formulaProduksi/:id", middleware.PermissionMiddleware("formula:delete"), handler.DeleteFormulaProduksi)

	//Barang Satuan Routes
	r.GET("/barangSatuan", middleware.PermissionMiddleware("satuan:read"), handler.ListBarangSatuan)
	r.GET("/barangSatuan/:id", middleware.PermissionMiddleware("satuan:read"), handler.GetBarangSatuanByID)
	r.POST("/barangSatuan", middleware.PermissionMiddleware("satuan:create"), handler.AddBarangSatuan)
	r.PUT("/barangSatuan/:id", middleware.PermissionMiddleware("satuan:update"), handler.UpdateBarangSatuan)
	r.DELETE("/barangSatuan/:id", middleware.PermissionMiddleware("satuan:delete"), handler.DeleteBarangSatuan)

	r.GET("/barangSatuanTurunan", middleware.PermissionMiddleware("satuan:read"), handler.ListSatuanTurunan)
	r.GET("/barangSatuanTurunan/:id", middleware.PermissionMiddleware("satuan:read"), handler.GetSatuanTurunanByID)
	r.POST("/barangSatuanTurunan", middleware.PermissionMiddleware("satuan:create"), handler.AddSatuanTurunan)
	r.PUT("/barangSatuanTurunan/:id", middleware.PermissionMiddleware("satuan:update"), handler.UpdateSatuanTurunan)
	r.DELETE("/barangSatuanTurunan/:id", middleware.PermissionMiddleware("satuan:delete"), handler.DeleteSatuanTurunan)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
