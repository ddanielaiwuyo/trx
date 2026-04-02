package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/persona-mp3/api"
	"github.com/persona-mp3/internal"
)

func main() {
	router := gin.Default()
	router.Use(internal.SetCORS())

	router.LoadHTMLGlob("frontend/templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/upload", api.UploadBankStatement)
	router.Run("localhost:8080")
}
