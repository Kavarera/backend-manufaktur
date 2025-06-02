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

	r.GET("/users", handler.AllUserList)
	r.GET("/users/:username", handler.UserList)
	r.DELETE("/users/:username", handler.UserDelete)

	// Public routes
	r.POST("/login", handler.Login)

	r.POST("/barangMentah", handler.AddMentah)
	r.GET("/barangMentah", handler.ListMentah)
	r.PUT("/barangMentah/:id", handler.UpdateMentah)
	r.DELETE("/barangMentah/:id", handler.DeleteMentah)

	r.POST("/perintahKerja", handler.AddPerintahKerja)
	r.GET("/perintahKerja", handler.ListPerintahKerja)
	r.PUT("/perintahKerja/:id", handler.UpdatePerintahKerja)
	r.DELETE("/perintahKerja/:id", handler.DeletePerintahKerja)
	r.POST("/perintahKerja/:id/upload-document", handler.UploadDocumentForPerintahKerja)
	r.GET("/perintahKerja/:id/download-document", handler.DownloadDocument)

	r.PUT("/updatePengerjaan/:id", handler.UpdateProsesPengerjaan)

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.RoleBasedAuth([]string{"SuperAdmin"}))
	{
		authGroup.POST("/register", handler.Register)
		authGroup.GET("/users", handler.AllUserList)
		authGroup.GET("/users/:username", handler.UserList)
		authGroup.DELETE("/users/:username", handler.UserDelete)
	}

	manageGroup := r.Group("/auth")
	manageGroup.Use(middleware.RoleBasedAuth([]string{"BarangManagement"}))
	{
		r.GET("/barangProduksi", handler.ListBarangProduksi)
		r.GET("/barangProduksi/:id", handler.GetBarangProduksiByID)
		r.POST("/barangProduksi", handler.AddBarangProduksi)
		r.PUT("/barangProduksi/:id", handler.UpdateBarangProduksi)
		r.DELETE("/barangProduksi/:id", handler.DeleteBarangProduksi)

		r.GET("/gudang", handler.ListGudang)
		r.GET("/gudang/:id", handler.GetGudangByID)
		r.POST("/gudang", handler.AddGudang)
		r.PUT("/gudang/:id", handler.UpdateGudang)
		r.DELETE("/gudang/:id", handler.DeleteGudang)
	}

	RPGroup := r.Group("/auth")
	RPGroup.Use(middleware.RoleBasedAuth([]string{"RencanaProduksi"}))
	{
		r.GET("/rencanaProduksi", handler.ListRencanaProduksi)
		r.GET("/rencanaProduksi/:id", handler.GetRencanaProduksiByID)
		r.POST("/rencanaProduksi", handler.AddRencanaProduksi)
		r.PUT("/rencanaProduksi/:id", handler.UpdateRencanaProduksi)
		r.DELETE("/rencanaProduksi/:id", handler.DeleteRencanaProduksi)
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
