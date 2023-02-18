package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"
	host := "localhost"
	route := fmt.Sprintf("%s:%s", host, port)

	router := gin.Default()
	router.GET("/", bookConverterMain)

	router.Run(route)
}

func bookConverterMain(c *gin.Context) {
	c.String(http.StatusOK, "Let's convert some books 📚")
}
