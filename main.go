package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// "http://localhost:3000"
type Bookmark struct {
	ID    int    `json:"id"`
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required,url"`
}

var bookmarks = []Bookmark{}
var nextID = 1

func main() {
	router := gin.Default()
	
	// Home Page Shows some cool UI.	
	router.LoadHTMLFiles("templates/index.html")
	router.Static("/static","./static")
	router.GET("/",func(c *gin.Context){
		c.HTML(200,"index.html",gin.H{
			"title":"Bookmarks API",
		})
	})

	// Read all bookmarks
	router.GET("/bookmarks", func(c *gin.Context) {
		c.JSON(200, bookmarks)
	})

	// Read one Bookmark
	router.GET("/bookmarks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for _, b := range bookmarks {
			if b.ID == id {
				c.JSON(200, b)
				return
			}
		}

		c.JSON(404, gin.H{"error": "Id does not exist"})
	})
	// Update a bookmark
	router.PUT("/bookmarks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		var updated Bookmark

		if err := c.ShouldBindJSON(&updated); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		for i, b := range bookmarks {
			if b.ID == id {
				updated.ID = id
				bookmarks[i] = updated
				c.JSON(200, updated)
				return
			}
		}

		c.JSON(404, gin.H{"error": "bookmark not found"})

	})

	// Delete a Bookmark
	router.DELETE("/bookmarks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		for i, b := range bookmarks {
			if b.ID == id {
				bookmarks = append(bookmarks[:i], bookmarks[i+1:]...)
			}
		}
	})

	// Create a Bookmark
	router.POST("/bookmarks", func(c *gin.Context) {
		var newBookmark Bookmark

		// JSON Binding
		if err := c.ShouldBindJSON(&newBookmark); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		newBookmark.ID = nextID
		nextID++

		bookmarks = append(bookmarks, newBookmark)
		c.JSON(201, newBookmark)
	})

	err := router.Run(":8080")
	if err != nil {
		println(" The Router encountered an error while running \n", err)
	}
}
