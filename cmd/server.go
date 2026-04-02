package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/persona-mp3/api"
	"github.com/persona-mp3/internal"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func main() {
	router := gin.Default()
	router.Use(internal.SetCORS())

	router.LoadHTMLGlob("frontend/templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/albums", addAlbum)
	router.POST("/upload", api.UploadBankStatement)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumById)
	router.Run("localhost:8080")
}
