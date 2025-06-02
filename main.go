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
	}

	PRGroup := r.Group("/auth")
	PRGroup.Use(middleware.RoleBasedAuth([]string{"PrintahKerja"}))
	{
		PRGroup.POST("/perintahKerja", handler.AddPerintahKerja)
		PRGroup.GET("/perintahKerja", handler.ListPerintahKerja)
		PRGroup.PUT("/perintahKerja/:id", handler.UpdatePerintahKerja)
		PRGroup.POST("/perintahKerja/:id/upload-document", handler.UploadDocumentForPerintahKerja)
		PRGroup.GET("/perintahKerja/:id/download-document", handler.DownloadDocument)
		PRGroup.PUT("/updatePengerjaan/:id", handler.UpdateProsesPengerjaan)
	}

	HPRGroup := r.Group("/auth")
	HPRGroup.Use(middleware.RoleBasedAuth([]string{"PrintahKerja"}))
	{
		PRGroup.DELETE("/perintahKerja/:id", handler.DeletePerintahKerja)
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
