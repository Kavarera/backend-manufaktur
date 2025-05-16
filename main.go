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
