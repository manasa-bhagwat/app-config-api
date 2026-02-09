package main

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	router.POST("/pages", CreatePage)
	router.GET("/pages", GetPages)
	router.GET("/pages/:id", GetPageByID)
	router.PUT("/pages/:id", UpdatePage)
	router.DELETE("/pages/:id", DeletePage)
}
