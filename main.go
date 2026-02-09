package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	// connect to DB
	InitDB()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	SetupRoutes(router)

	router.Run(":8080")
}
