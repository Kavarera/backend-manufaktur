package main

import (
	"ProjekRapli2/API/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	r.POST("/login", handler.login)

	registerGroup := r.Group("/auth")
	registerGroup.Use()

}
