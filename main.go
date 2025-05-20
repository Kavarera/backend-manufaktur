package main

import (
	"log"
	"manufacture_API/db"
	"manufacture_API/handler"
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
	r.POST("/register", handler.Register)
	r.GET("/users", handler.AllUserList)
	r.GET("/users/:id", handler.UserList)
	r.DELETE("/users/:id", handler.UserDelete)

	r.POST("/barangMentah", handler.AddMentah)
	r.GET("/barangMentah", handler.ListMentah)
	r.PUT("/barangMentah/:id", handler.UpdateMentah)
	r.DELETE("/barangMentah/:id", handler.DeleteMentah)

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

	r.GET("/rencanaProduksi", handler.ListRencanaProduksi)
	r.GET("/rencanaProduksi/:id", handler.GetRencanaProduksiByID)
	r.POST("/rencanaProduksi", handler.AddRencanaProduksi)
	r.PUT("/rencanaProduksi/:id", handler.UpdateRencanaProduksi)
	r.DELETE("/rencanaProduksi/:id", handler.DeleteRencanaProduksi)

	// // Protected routes group with RoleBasedAuth middleware
	// authGroup := r.Group("/auth")
	// authGroup.Use(middlewares.RoleBasedAuth([]string{"Super Admin"}))
	// {
	// 	// authGroup.GET("/register", handler.Register) // Assuming exported function Register
	// }

	// changeGroup := r.Group("/auth")
	// changeGroup.Use(middlewares.RoleBasedAuth([]string{"User"}))
	// {
	// 	// changeGroup.GET("/register", handler.Register) // Assuming exported function Register
	// }

	// Run server on port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
