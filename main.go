package main

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := "4444"
	host := "localhost"
	route := fmt.Sprintf("%s:%s", host, port)

	router := gin.Default()

	// Enable CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	router.Use(cors.New(config))

	// Define endpoints
	router.GET("/", greet)
	router.GET("/convert/book", convertBook)
	router.GET("/convert/author", convertAuthor)
	router.StaticFile("author.jpg", "./static/author.jpg")

	router.Run(route)
}

func greet(c *gin.Context) {
	c.String(http.StatusOK, "Let's convert some books 📚")
}

func convertBook(c *gin.Context) {
	query := c.Request.URL.Query()
	url2conv := query.Get("link")
	valid := govalidator.IsURL(url2conv)

	if !valid {
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
		return
	}

	book, err := scrapeBook(url2conv)

	if err != nil || book == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSONP(http.StatusOK, book)
}

func convertAuthor(c *gin.Context) {
	query := c.Request.URL.Query()
	url2conv := query.Get("link")
	valid := govalidator.IsURL(url2conv)

	if !valid {
		c.IndentedJSON(http.StatusBadRequest, gin.H{})
		return
	}

	author, err := scrapeAuthor(url2conv)

	if err != nil || author == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSONP(http.StatusOK, author)
}
