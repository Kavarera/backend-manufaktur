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

	// Create Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Public routes
	r.POST("/login", handler.Login)
	r.GET("/jadwalProduksi", handler.ListRencanaProduksi)
	r.POST("/rencanaProduksi", handler.AddRencanaProduksi)
	r.GET("/rencanaProduksi", handler.ListRencanaProduksi)
	r.PUT("/rencanaProduksi/:id", handler.UpdateRencanaProduksi)
	r.DELETE("/rencanaProduksi/:id", handler.DeleteRencanaProduksi)

	authGroup := r.Group("/admin")
	authGroup.Use(middleware.RoleBasedAuth([]string{"SuperAdmin"}))
	{
		authGroup.POST("/register", handler.Register)
		authGroup.GET("/users", handler.AllUserList)
		authGroup.GET("/users/:username", handler.UserList)
		authGroup.DELETE("/users/:username", handler.UserDelete)

		authGroup.GET("/barangProduksi", handler.ListBarangProduksi)
		authGroup.GET("/barangProduksi/:id", handler.GetBarangProduksiByID)
		authGroup.POST("/barangProduksi", handler.AddBarangProduksi)
		authGroup.PUT("/barangProduksi/:id", handler.UpdateBarangProduksi)
		authGroup.DELETE("/barangProduksi/:id", handler.DeleteBarangProduksi)

		authGroup.GET("/gudang", handler.ListGudang)
		authGroup.GET("/gudang/:id", handler.GetGudangByID)
		authGroup.POST("/gudang", handler.AddGudang)
		authGroup.PUT("/gudang/:id", handler.UpdateGudang)
		authGroup.DELETE("/gudang/:id", handler.DeleteGudang)

		authGroup.POST("/barangMentah", handler.AddMentah)
		authGroup.GET("/barangMentah", handler.ListMentah)
		authGroup.PUT("/barangMentah/:id", handler.UpdateMentah)
		authGroup.DELETE("/barangMentah/:id", handler.DeleteMentah)

		authGroup.GET("/rencanaProduksi", handler.ListRencanaProduksi)
		authGroup.GET("/rencanaProduksi/:id", handler.GetRencanaProduksiByID)
		authGroup.POST("/rencanaProduksi", handler.AddRencanaProduksi)
		authGroup.PUT("/rencanaProduksi/:id", handler.UpdateRencanaProduksi)
		authGroup.DELETE("/rencanaProduksi/:id", handler.DeleteRencanaProduksi)
		authGroup.GET("/jadwalProduksi", handler.ListRencanaProduksi)

		authGroup.POST("/perintahKerja", handler.AddPerintahKerja)
		authGroup.GET("/perintahKerja", handler.ListPerintahKerja)
		authGroup.PUT("/perintahKerja/:id", handler.UpdatePerintahKerja)
		authGroup.POST("/perintahKerja/:id/upload-document", handler.UploadDocumentForPerintahKerja)
		authGroup.GET("/perintahKerja/:id/download-document", handler.DownloadDocument)
		authGroup.PUT("/updatePengerjaan/:id", handler.UpdateProsesPengerjaan)
		authGroup.DELETE("/perintahKerja/:id", handler.DeletePerintahKerja)

		authGroup.POST("/pengambilanBarangBaku", handler.AddPengambilanBarangBaku)
		authGroup.GET("/pengambilanBarangBaku", handler.GetPengambilanBarangBaku)
		authGroup.PUT("/pengambilanBarangBaku/:id", handler.UpdatePengambilanBarangBaku)
		authGroup.DELETE("/pengambilanBarangBaku/:id", handler.DeletePengambilanBarangBaku)
	}

	manageGroup := r.Group("/auth")
	manageGroup.Use(middleware.RoleBasedAuth([]string{"BarangManagement"}))
	{
		manageGroup.GET("/barangProduksi", handler.ListBarangProduksi)
		manageGroup.GET("/barangProduksi/:id", handler.GetBarangProduksiByID)
		manageGroup.POST("/barangProduksi", handler.AddBarangProduksi)
		manageGroup.PUT("/barangProduksi/:id", handler.UpdateBarangProduksi)
		manageGroup.DELETE("/barangProduksi/:id", handler.DeleteBarangProduksi)

		manageGroup.GET("/gudang", handler.ListGudang)
		manageGroup.GET("/gudang/:id", handler.GetGudangByID)
		manageGroup.POST("/gudang", handler.AddGudang)
		manageGroup.PUT("/gudang/:id", handler.UpdateGudang)
		manageGroup.DELETE("/gudang/:id", handler.DeleteGudang)

		manageGroup.POST("/barangMentah", handler.AddMentah)
		manageGroup.GET("/barangMentah", handler.ListMentah)
		manageGroup.PUT("/barangMentah/:id", handler.UpdateMentah)
		manageGroup.DELETE("/barangMentah/:id", handler.DeleteMentah)
	}

	RPGroup := r.Group("/auth")
	RPGroup.Use(middleware.RoleBasedAuth([]string{"RencanaProduksi"}))
	{
		RPGroup.GET("/rencanaProduksi", handler.ListRencanaProduksi)
		RPGroup.GET("/rencanaProduksi/:id", handler.GetRencanaProduksiByID)
		RPGroup.POST("/rencanaProduksi", handler.AddRencanaProduksi)
		RPGroup.PUT("/rencanaProduksi/:id", handler.UpdateRencanaProduksi)
		RPGroup.DELETE("/rencanaProduksi/:id", handler.DeleteRencanaProduksi)
		RPGroup.GET("/jadwalProduksi", handler.ListRencanaProduksi)
	}

	PRGroup := r.Group("/auth")
	PRGroup.Use(middleware.RoleBasedAuth([]string{"PerintahKerja"}))
	{
		PRGroup.POST("/perintahKerja", handler.AddPerintahKerja)
		PRGroup.GET("/perintahKerja", handler.ListPerintahKerja)
		PRGroup.PUT("/perintahKerja/:id", handler.UpdatePerintahKerja)
		PRGroup.POST("/perintahKerja/:id/upload-document", handler.UploadDocumentForPerintahKerja)
		PRGroup.GET("/perintahKerja/:id/download-document", handler.DownloadDocument)
		PRGroup.PUT("/updatePengerjaan/:id", handler.UpdateProsesPengerjaan)
	}

	HPRGroup := r.Group("/auth")
	HPRGroup.Use(middleware.RoleBasedAuth([]string{"HapusPerintahKerja"}))
	{
		HPRGroup.DELETE("/perintahKerja/:id", handler.DeletePerintahKerja)
	}

	PBKGroup := r.Group("/auth")
	PBKGroup.Use(middleware.RoleBasedAuth([]string{"PengambilanBarangBaku"}))
	{
		PBKGroup.POST("/pengambilanBarangBaku", handler.AddPengambilanBarangBaku)
		PBKGroup.GET("/pengambilanBarangBaku", handler.GetPengambilanBarangBaku)
		PBKGroup.PUT("/pengambilanBarangBaku/:id", handler.UpdatePengambilanBarangBaku)
		PBKGroup.DELETE("/pengambilanBarangBaku/:id", handler.DeletePengambilanBarangBaku)
	}

	PBJGroup := r.Group("/auth")
	PBJGroup.Use(middleware.RoleBasedAuth([]string{"PengambilanBarangJadi"}))
	{

	}

	// Run server on port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
