package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Basic ping route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Great 🏓",
		})
	})

	// Hello route
	router.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("world")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + name + " 👋",
		})
	})

	router.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080
}
